var app = new Vue({
  el: "#app",
  data: {
    version: "v1"
  },
  computed: {
    host: function() {
      const protocol = location.protocol
      const slashes = protocol.concat("//")
      return slashes.concat(window.location.hostname)
    },
    authSlack: function() {
      return `${this.host}/${this.version}/auth/slack`
    }
  }
})
