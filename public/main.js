var app = new Vue({
  el: '#app',
  data: {
    version: 'v1'
  },
  computed: {
    host: function() {
      const protocol = location.protocol;
      const slashes = protocol.concat('//');
      return slashes.concat(window.location.hostname);
    },
    port: function() {
      return location.port;
    },
    authSlack: function() {
      return `${this.host}:${this.port}/${this.version}/auth/slack`;
    }
  },
  methods: {
    auth() {
      popupCenter(this.authSlack, 'Authentication', '800', '600');
      return false;
    }
  }
});

function popupCenter(url, title, w, h) {
  // Fixes dual-screen position                         Most browsers      Firefox
  var dualScreenLeft = window.screenLeft != undefined ? window.screenLeft : window.screenX;
  var dualScreenTop = window.screenTop != undefined ? window.screenTop : window.screenY;

  var width = window.innerWidth ? window.innerWidth : document.documentElement.clientWidth ? document.documentElement.clientWidth : screen.width;
  var height = window.innerHeight ? window.innerHeight : document.documentElement.clientHeight ? document.documentElement.clientHeight : screen.height;

  var systemZoom = width / window.screen.availWidth;
  var left = (width - w) / 2 / systemZoom + dualScreenLeft;
  var top = (height - h) / 2 / systemZoom + dualScreenTop;
  var newWindow = window.open(url, title, 'scrollbars=yes, width=' + w / systemZoom + ', height=' + h / systemZoom + ', top=' + top + ', left=' + left);

  // Puts focus on the newWindow
  if (window.focus) newWindow.focus();
}
