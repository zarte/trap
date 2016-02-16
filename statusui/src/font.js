WebFontConfig = {
    google: { families: [ 'Oswald:400,700:latin', 'Roboto+Mono:400:latin' ] }
};

$(window).ready(function() {
    var wf = document.createElement('script'),
        s = document.getElementsByTagName('script')[0];;

    wf.src = ('https:' == document.location.protocol ? 'https' : 'http') +
    '://ajax.googleapis.com/ajax/libs/webfont/1/webfont.js';

    wf.type = 'text/javascript';
    wf.async = 'true';

    s.parentNode.insertBefore(wf, s);
});