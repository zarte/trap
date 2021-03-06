/*
 * Trap
 * An anti-pryer server for better privacy
 *
 * This file is a part of Trap project
 *
 * Copyright 2016 Rain Lee <raincious@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

 var initApp = function() {
    new Vue({
        el: '#p-body',
        data: {
            status: {
                nav: {
                    record: 'clients'
                },
                verifiy: {
                    running: false
                },
                session: {
                    loggedIn: false,
                    authID: '',
                    toString: function() {
                        return this.authID;
                    },
                    permissions: {}
                }
            },
            inputs: {
                password: ''
            },
            records: {
                source: {
                    clients: {
                        data: [],
                        listOrder: [],
                        clientMap: {},
                        loaded: false,
                        isLoaded: function(vueObj) {
                            return vueObj.records.source.clients.loaded;
                        },
                        fetcher: function(vueObj, finished) {
                            vueObj.requestJson('GET', '/api/clients', {}, function(data) {
                                vueObj.clientList = data;
                                vueObj.records.source.clients.loaded = true;
                            }, function(jqXHR, textStatus, errorThrown) {
                                switch (jqXHR.status) {
                                    case 403:
                                        vueObj.auth = null;
                                        break;
                                }
                            }, finished);
                        }
                    },
                    sessions: {
                        data: [],
                        loaded: false,
                        isLoaded: function(vueObj) {
                            return vueObj.records.source.sessions.loaded;
                        },
                        fetcher: function(vueObj, finished) {
                            vueObj.requestJson('GET', '/api/sessions', {}, function(data) {
                                var d = [];

                                for (var i in data) {
                                    d.push({
                                        IP:             data[i].IP,
                                        Created:        new Date(Date.parse(data[i].Created)),
                                        LastSeen:       new Date(Date.parse(data[i].LastSeen)),
                                        Expire:         data[i].Expire / 1000000,
                                        Permissions:    data[i].Permissions
                                    });
                                }

                                d.sort(function(a, b) {
                                    return b.LastSeen.getTime() - a.LastSeen.getTime();
                                });

                                vueObj.records.source.sessions.data = d;
                                vueObj.records.source.sessions.loaded = true;
                            }, function(jqXHR, textStatus, errorThrown) {
                                switch (jqXHR.status) {
                                    case 403:
                                        vueObj.auth = null;
                                        break;
                                }
                            }, finished);
                        }
                    },
                    sync: {
                        data: {
                            server: {
                                ip: "",
                                port: 0
                            },
                            partners: [],
                            indirects: [],
                        },
                        loaded: false,
                        isLoaded: function(vueObj) {
                            return vueObj.records.source.sync.loaded;
                        },
                        fetcher: function(vueObj, finished) {
                            vueObj.requestJson('GET', '/api/sync', {}, function(data) {
                                var partners = [];
                                var indirects = [];

                                for (var i in data.Nodes) {
                                    partners.push({
                                        IP: data.Nodes[i].Address.IP,
                                        Port: data.Nodes[i].Address.Port,
                                        Delay: data.Nodes[i].Delay,
                                        Rx: data.Nodes[i].Stats.RX,
                                        Tx: data.Nodes[i].Stats.TX,
                                        Connected: data.Nodes[i].Connected,
                                        Type: "Initiative",
                                        Weight: data.Nodes[i].Stats.RX + data.Nodes[i].Stats.TX
                                    });

                                    for (var j in data.Nodes[i].Partner) {
                                        indirects.push({
                                            IP: data.Nodes[i].Partner[j].IP,
                                            Port: data.Nodes[i].Partner[j].Port,
                                            Via: {
                                                IP: data.Nodes[i].Address.IP,
                                                Port: data.Nodes[i].Address.Port
                                            }
                                        });
                                    }
                                }

                                for (var i in data.Server.Clients) {
                                    partners.push({
                                        IP: data.Server.Clients[i].Remote.IP,
                                        Port: data.Server.Clients[i].Remote.Port,
                                        Delay: null,
                                        Rx: data.Server.Clients[i].Stats.RX,
                                        Tx: data.Server.Clients[i].Stats.TX,
                                        Connected: true,
                                        Type: "Passive",
                                        Weight: data.Server.Clients[i].Stats.RX + data.Server.Clients[i].Stats.TX
                                    });
                                }

                                partners.sort(function(a, b) {
                                    return (b.Connected - a.Connected
                                        || b.Weight - a.Weight
                                        || a.IP.localeCompare(b.IP)
                                        || a.Port - b.Port);
                                });

                                indirects.sort(function(a, b) {
                                    return (a.Via.IP.localeCompare(b.Via.IP)
                                        || a.Via.Port - b.Via.Port
                                        || a.IP.localeCompare(b.IP)
                                        || a.Port - b.Port);
                                });

                                vueObj.records.source.sync.data = {
                                    server: {
                                        ip: data.Server.IP,
                                        port: data.Server.Port
                                    },
                                    partners: partners,
                                    indirects: indirects
                                };
                                vueObj.records.source.sync.loaded = true;
                            }, function(jqXHR, textStatus, errorThrown) {
                                switch (jqXHR.status) {
                                    case 403:
                                        vueObj.auth = null;
                                        break;
                                }
                            }, finished);
                        }
                    },
                    logs: {
                        data: [],
                        loaded: false,
                        isLoaded: function(vueObj) {
                            return vueObj.records.source.logs.loaded;
                        },
                        fetcher: function(vueObj, finished) {
                            vueObj.requestJson('GET', '/api/logs', {}, function(data) {
                                var d = [];

                                for (var i in data) {
                                    d.push({
                                        Time:       new Date(Date.parse(data[i].Time)),
                                        Type:       data[i].Type,
                                        Context:    data[i].Context,
                                        Message:    data[i].Message
                                    });
                                }

                                d.sort(function(a, b) {
                                    return b.Time.getTime() - a.Time.getTime();
                                });

                                vueObj.records.source.logs.data = d;
                                vueObj.records.source.logs.loaded = true;
                            }, function(jqXHR, textStatus, errorThrown) {
                                switch (jqXHR.status) {
                                    case 403:
                                        vueObj.auth = null;
                                        break;
                                }
                            }, finished);
                        }
                    }
                },
                fetch: function(vueObj, showLoading, finished) {
                    var finishedCB = finished ? finished : function() {};

                    if (showLoading || !vueObj.records.source[vueObj.currentRecordType].isLoaded(vueObj)) {
                        vueObj.records.showingLoading = true;
                    }

                    vueObj.records.source[vueObj.currentRecordType].fetcher(vueObj, function() {
                        vueObj.records.showingLoading = false;

                        finishedCB();
                    });
                },
                showingLoading: false,
                queuer: new jQueue(function(vueObj) {
                    vueObj.records.fetch(vueObj);
                }, 'Loop', 30000)
            },
            charts: {
                status: {
                    synced: false,
                    totalInbound: 0,
                    totalMarked: 0,
                    totalHit: 0,
                    totalClients: 0,
                    uptime: 0
                },
                syncer: function(vueObj) {
                    vueObj.requestJson('GET', '/api/status', {}, function(data) {
                        var parseDistributionData = function(distribution, maxItems) {
                                var totalHit = 0,
                                    portAccesses = {},
                                    portPercents = [],
                                    result = [];

                                for (var i in distribution) {
                                    portAccesses[distribution[i].Port + ':' + distribution[i].Type] = distribution[i];

                                    totalHit += distribution[i].Hit;
                                }

                                for (var i in portAccesses) {
                                    portPercents.push({
                                        Percent: (portAccesses[i].Hit / totalHit) * 100,
                                        Type: portAccesses[i].Type + ' ' + portAccesses[i].Port,
                                        Port: portAccesses[i].Port
                                    });
                                }

                                portPercents.sort(function(a, b) {
                                    return b.Percent - a.Percent || a.Type.localeCompare(b.Type);
                                });

                                for (var i in portPercents) {
                                    if (result.length > maxItems) {
                                        result[maxItems].Percent += portPercents[i].Percent;
                                        result[maxItems].Type = 'Rest';
                                        result[maxItems].Port = 'Rest';

                                        continue;
                                    }

                                    result.push(portPercents[i]);
                                }

                                return result;
                            },
                            parseHistoryData = function(history, length, curRound) {
                                var maxRound = 0,
                                    rounds = {},
                                    result = [];

                                for (var i in history) {
                                    rounds[history[i].Hours] = history[i];

                                    if (history[i].Hours > maxRound) {
                                        maxRound = history[i].Hours;
                                    }
                                }

                                if (curRound > maxRound) {
                                    maxRound = curRound;
                                }

                                for (var i = maxRound; i > 0; --i) {
                                    if (result.length >= length) {
                                        break;
                                    }

                                    if (typeof rounds[i] !== 'object') {
                                        result.push({
                                            Marked: 0,
                                            Inbound: 0,
                                            Hit: 0,
                                            Hours: i
                                        });

                                        continue;
                                    }

                                    result.push(rounds[i]);
                                }

                                for (var i = result.length; i < length; ++i) {
                                    result.push({
                                        Marked: 0,
                                        Inbound: 0,
                                        Hit: 0,
                                        Hours: 0
                                    });
                                }

                                return result;
                            },
                            history = {
                                mk: [],
                                ib: [],
                                ht: [],
                                lb: []
                            },
                            distribution = {
                                d: [],
                                n: []
                            },
                            currentUpHour = data.Uptime / 3600000000000,
                            parsedHistory = parseHistoryData(
                                data.History,
                                12,
                                Math.ceil(currentUpHour)
                            ),
                            parsedDistribution = (data.Distribution.length > 0 ? parseDistributionData(
                                data.Distribution,
                                5
                            ) : []);

                        vueObj.charts.status.synced         =   true;

                        vueObj.charts.status.totalInbound   =   data.TotalInbound;
                        vueObj.charts.status.totalMarked    =   data.TotalMarked;
                        vueObj.charts.status.totalHit       =   data.TotalHit;
                        vueObj.charts.status.totalClients   =   data.TotalClients;
                        vueObj.charts.status.uptime         =   Math.round(currentUpHour);

                        for (var i in parsedHistory) {
                            history.mk.push(parsedHistory[i].Marked);
                            history.ib.push(parsedHistory[i].Inbound);
                            history.ht.push(parsedHistory[i].Hit);
                            history.lb.push(parsedHistory[i].Hours);
                        }

                        vueObj.charts.history.update({
                            labels: history.lb,
                            series: [
                                history.ht,
                                history.ib,
                                history.mk
                            ]
                        });

                        for (var i in parsedDistribution) {
                            distribution.n.push(parsedDistribution[i].Type);
                            distribution.d.push(parsedDistribution[i].Percent);
                        }

                        vueObj.charts.accesses.update({
                            labels: (distribution.n.length < 1 ? [' '] : distribution.n),
                            series: (distribution.d.length < 1 ? [100] : distribution.d)
                        });
                    }, function(jqXHR, textStatus, errorThrown) {
                        switch (jqXHR.status) {
                            case 403:
                                vueObj.auth = null;
                                break;
                        }
                    });
                },
                queuer: new jQueue(function(vueObj) {
                    vueObj.charts.syncer(vueObj);
                }, 'Loop', 20000),
                history: new Chartist.Line(
                    '#status-stats-history-chart',
                    {
                        labels: [12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1],
                        series: [
                            [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
                        ]
                    },
                    {
                        showArea: true,
                        showLine: false,
                        showPoint: false,
                        fullWidth: true,
                        low: 0,
                        axisX: {
                            showLabel: true,
                            showGrid: false
                        }
                    }
                ),
                accesses: new Chartist.Pie(
                    '#status-stats-distro-chart',
                    {
                        labels: [
                            ''
                        ],
                        series: [
                            100
                        ]
                    },
                    {
                        total: 100,
                        donut: true,
                        donutWidth: 15,
                        startAngle: 0,
                        showLabel: true
                    }
                )
            }
        },
        computed: {
            clientList: {
                get: function() {
                    this.records.source.clients.data = [];

                    if (this.records.source.clients.listOrder.length <= 0) {
                        return this.records.source.clients.data;
                    }

                    for (var i in this.records.source.clients.listOrder) {
                        this.records.source.clients.data.push(
                            this.records.source.clients.clientMap[
                                this.records.source.clients.listOrder[i].Key
                            ]);
                    }

                    return this.records.source.clients.data;
                },
                set: function(clientList) {
                    var vueObj = this,
                        updateableAttributes = ['LastSeen', 'Count', 'Records', 'Marked'],
                        newClientKeys = {},
                        parseClientData = function(clientData) {
                            return {
                                Address:        clientList[i].Address,
                                LastSeen:       new Date(Date.parse(clientList[i].LastSeen)),
                                FirstSeen:      new Date(Date.parse(clientList[i].FirstSeen)),
                                Count:          clientList[i].Count,
                                Records:        clientList[i].Records,
                                Marked:         clientList[i].Marked,
                                RecordData:     [],
                                Expended:       false,
                                Expend:         function(index) {
                                    if (this.Expended) {
                                        this.Expended = false;

                                        $('#status-marked-client-data-' + index).slideUp(500);

                                        return;
                                    }

                                    this.RecordData = vueObj.parseClientRecords(this.Records);
                                    this.Expended = true;

                                    $('#status-marked-client-data-' + index).slideDown(500);
                                },
                                Deleting:       false,
                                Delete:         function(index) {
                                    if (this.Deleting) {
                                        return;
                                    }

                                    if (!confirm('Do you want to remove and unmark inbound client \'' + this.Address + '\'?')) {
                                        return;
                                    }

                                    this.Deleting = true;
                                    vueObj.records.queuer.reset();

                                    vueObj.requestJson(
                                        'DELETE',
                                        '/api/client?client=' + encodeURIComponent(this.Address),
                                        {},
                                        function(data) {
                                            this.Deleting = false;

                                            if (!data.Result) {
                                                alert('Error happened :(');

                                                return;
                                            }

                                            $('#status-marked-client-' + index).slideUp(
                                                500,
                                                function() {
                                                    // Manually fetch updated client data
                                                    vueObj.records.fetch(vueObj, false, function() {
                                                        // Restart data sync after data load
                                                        vueObj.records.queuer.run(vueObj);
                                                    });
                                                }
                                            );
                                        },
                                        function(jqXHR, textStatus, errorThrown) {
                                            switch (jqXHR.status) {
                                                case 403:
                                                case 400:
                                                case 401:
                                                    vueObj.auth = null;
                                                    break;
                                            }

                                            if (typeof jqXHR.responseJSON !== 'object') {
                                                alert('Error happened :(');

                                                return;
                                            }

                                            alert(jqXHR.responseJSON.Error);
                                        }
                                    );
                                }
                            };
                        },
                        eClient = null;

                    this.records.source.clients.listOrder = [];

                    for (var i in clientList) {
                        newClientKeys[clientList[i].Address] = true;

                        eClient = parseClientData(clientList[i]);

                        this.records.source.clients.listOrder.push({
                            Key:        eClient.Address,
                            LastSeen:   eClient.LastSeen,
                            Count:      eClient.Count
                        });

                        if (typeof this.records.source.clients.clientMap[eClient.Address] !== 'object') {
                            // If this is a new client, add it to the map
                            this.records.source.clients.clientMap[eClient.Address] = eClient;

                            continue;
                        }

                        // If this is a existing client, update it
                        for (var a in updateableAttributes) {
                            this.records.source.clients.clientMap[eClient.Address][updateableAttributes[a]] =
                                eClient[updateableAttributes[a]];
                        }

                        // If current client record is expanded, update clientRecord as well
                        if (this.records.source.clients.clientMap[eClient.Address].Expended) {
                            this.records.source.clients.clientMap[eClient.Address].RecordData =
                                this.parseClientRecords(eClient.Records);
                        }
                    }

                    this.records.source.clients.listOrder.sort(function(a, b) {
                        return ((b.Count - a.Count)
                            || (b.LastSeen.getTime() - a.LastSeen.getTime())
                            || a.Key.localeCompare(b.Key));
                    });

                    // Update clients count
                    this.charts.status.totalClients =
                        this.records.source.clients.listOrder.length;

                    // Scan deleted clients from client map
                    for (var c in this.records.source.clients.clientMap) {
                        if (typeof newClientKeys[c] !== 'undefined') {
                            continue;
                        }

                        delete this.records.source.clients.clientMap[c];
                    }
                }
            },
            auth: {
                get: function() {
                    return this.status.session;
                },
                set: function (session) {
                    if (session !== null) {
                        $('#p-login').slideUp(1000);
                        $('#p-main').slideDown(1000);

                        // Reset charts
                        this.charts.history.update();
                        this.charts.accesses.update();

                        // Start stats data sync
                        this.charts.queuer.run(this);
                        this.records.queuer.run(this);

                        this.status.session = {
                            loggedIn: true,
                            authID: session.Token,
                            permissions: session.Permissions
                        }

                        return;
                    }

                    $('#p-main').slideUp(1000);
                    $('#p-login').slideDown(1000);

                    // End stats data sync
                    this.charts.queuer.reset();

                    // End record data sync
                    this.records.queuer.reset();

                    this.status.session = {
                        loggedIn: false,
                        authID: '',
                        permissions: {}
                    };

                    setTimeout(function() {
                        $('#status-login-password').focus();
                        $('#status-login-password').click();
                    }, 500);
                }
            },
            currentRecordType: {
                get: function() {
                    return this.status.nav.record;
                },
                set: function (recordType) {
                    var switching = false;

                    if (typeof this.records.source[recordType] !== 'object') {
                        return;
                    }

                    if (this.status.nav.record === recordType) {
                        return;
                    }

                    this.status.nav.record = recordType;

                    // Manually fetch data once
                    this.records.fetch(this);

                    // Reset the queuer for new nav type
                    this.records.queuer.reset();
                    this.records.queuer.run(this);
                }
            }
        },
        methods: {
            parseTime: function(dateTime) {
                var d = new Date(dateTime),
                    s = [   'Jan', 'Feb', 'Mar', 'Apr',
                            'May', 'Jun', 'Jul', 'Aug',
                            'Sep', 'Oct', 'Nov', 'Dec'];

                return {
                    Clock:  d.getHours() + ':' + d.getMinutes() + ':' + d.getSeconds(),
                    Date:   s[d.getMonth()] + ', ' + d.getDate(),
                    Year:   d.getFullYear()
                };
            },
            numFilter: function(num, formats) {
                var n = parseInt(num, 10), format = [0, formats[0]];

                for (var f in formats) {
                    if (n < f) {
                        break;
                    }

                    format = [f, formats[f]];
                }

                if (format[0] == 0) {
                    return format[1].replace(
                        '%N%',
                        Math.round(n).toLocaleString()
                    );
                }

                return format[1].replace(
                    '%N%',
                    Math.round(n / format[0]).toLocaleString()
                );
            },
            httpRequest: function(method, url, reqData, reqDataType, succ, fail, always) {
                var headers = {}, empty = {};
                var successCb = typeof succ === 'function' ? succ : function() {},
                    failCb = typeof fail === 'function' ? fail : function() {},
                    alwaysCb = typeof always === 'function' ? always : function() {},
                    rDType = typeof reqDataType === 'string' ? reqDataType : 'json'
                    rD = typeof reqData !== 'undefined' && reqData ? reqData : '';

                if (this.auth.authID !== '') {
                    headers['X-Trap-Token'] = this.auth.authID;
                }

                $.ajax({
                    url: url,
                    cache: false,
                    contentType: 'application/json; charset=utf-8',
                    dataType: rDType,
                    data: rD,
                    headers: headers,
                    method: method,
                    success: function(data) {
                        if (typeof data !== 'object') {
                            return;
                        }

                        successCb(data);
                    },
                    error: function(jqXHR, textStatus, errorThrown) {
                        if (typeof jqXHR !== 'object') {
                            return;
                        }

                        failCb(jqXHR, textStatus, errorThrown);
                    },
                    complete: function() {
                        alwaysCb();
                    }
                });
            },
            requestJson: function(method, url, data, succ, failed, always) {
                this.httpRequest(
                    method,
                    url,
                    JSON.stringify(data),
                    'json',
                    succ,
                    failed,
                    always
                );
            },
            subNavSwitch: function(nav) {
                this.currentRecordType = nav;
            },
            verifySession: function(passwd) {
                var self = this;

                if (passwd === '') {
                    return;
                }

                if (self.status.verifiy.running) {
                    return;
                }

                self.status.verifiy.running = true;

                this.requestJson('POST', '/api/auth', {
                    Password: passwd
                }, function(data) {
                    self.auth = {
                        Token:          data.Token,
                        Permissions:    data.Permissions
                    };

                    // Reload Charts
                    self.charts.syncer(self);

                    // Reload Records
                    self.records.fetch(self, true);
                }, function(jqXHR, textStatus, errorThrown) {
                    switch (jqXHR.status) {
                        case 403:
                        case 400:
                        case 401:
                            self.auth = null;
                            break;
                    }

                    if (typeof jqXHR.responseJSON !== 'object') {
                        alert('Error happened :(');

                        return;
                    }

                    alert(jqXHR.responseJSON.Error);
                }, function() {
                    self.status.verifiy.running = false;
                });
            },
            parseClientRecords: function(data) {
                var result = [];

                for (var i in data) {
                    result.unshift(this.parseClientRecord(data[i]));
                }

                return result;
            },
            parseClientRecord: function(data) {
                var date = this.parseTime(data.Time),
                    inbound = atob(data.Inbound),
                    outbound = atob(data.Outbound);

                return {
                    Request: inbound,
                    Respond: outbound,
                    Hitting: data.Hitting,
                    Time: date
                };
            }
        },
        filters: {
            trueFields: function(fs) {
                var result = [];

                for (var i in fs) {
                    if (!fs[i]) {
                        continue;
                    }

                    result.push(i);
                }

                if (result.length <= 0) {
                    return '<span>N/A</span>';
                }

                return '<em>' + result.join('</em><em>') + '</em>';
            },
            spiltedBy: function(val, spilter) {
                return '<em>' + val.split(spilter).join('</em><em>') + '</em>';
            },
            dayTime: function(ts) {
                var time = parseInt(ts.getTime(), 10),
                    ago = (+ new Date()) - time,
                    formats = {
                        0: 'just now',
                        1000: '%N% seconds ago',
                        60000: '%N% minutes ago',
                        3600000: '%N% hours ago',
                        86400000: '%N% days ago',
                        2592000000: '%N% months ago',
                        31104000000: '%N% years ago'
                    }

                return this.numFilter(ago, formats);
            },
            duration: function(d) {
                if (d <= 0) {
                    return 'n/a';
                }

                return this.numFilter(d, {
                    0: '%N%ns',
                    1000: '%N%μs',
                    1000000: '%N%ms',
                    1000000000: '%N%s',
                    60000000000: '%N%m',
                    360000000000: '%N%h'
                });
            },
            bytes: function(b) {
                if (b <= 0) {
                    return '0B';
                }

                return this.numFilter(b, {
                    0: '%N%B',
                    1024: '%N%KiB',
                    1048576: '%N%MiB',
                    1073741824: '%N%GiB'
                });
            },
            ASCIICode: function(str) {
                var codes = '';

                for (var i = 0; i < str.length; ++i) {
                    var tCode = str.charCodeAt(i).toString(16).toUpperCase();

                    switch (tCode.length) {
                        case 0:
                            tCode = '00';
                            break;

                        case 1:
                            tCode = '0' + tCode;
                            break;
                    }

                    codes += '<em>' + tCode + '</em>';
                }

                if (codes === '') {
                    return '';
                }

                return codes;
            },
            number: function(n) {
                return n.toLocaleString();
            },
            defaultVal: function(v, d) {
                if (v) {
                    return v;
                }

                return d;
            }
        }
    });
};
