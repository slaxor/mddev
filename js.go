package main

const JS = `(function() {
  var md = document.getElementById("markdown");
  var conn = new WebSocket("ws://" + location.host +
      "/ws?lastMod=" + (new Date().valueOf()));
  conn.onclose = function(evt) {
    md.innerHTML = 'Connection closed';
  }
  conn.onmessage = function(evt) {
    console.log('file updated');
    md.innerHTML = evt.data;
  }
})();
`
