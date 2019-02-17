const etc = require("etc-js");
const Core = require("./core");
const inherits     = require("inherits");
const EventEmitter = require("events").EventEmitter;

inherits(Webdaemon, EventEmitter);

function Webdaemon(file) {
  this.conn = new Core(file, this.onMessage.bind(this));
  this.query = {};

  EventEmitter.call(this);
}

Webdaemon.prototype.onMessage = function(event) {
  const $this = this;

  switch (event.type) {
    case "answer":
    if ($this.query[event.id]) {
      $this.query[event.id](event.data);
    }
    break;
  }
}

Webdaemon.prototype.q = function(method, body) {
  const $this = this;

  return new Promise(function(a, rj) {
    var id = new etc.UUID();

    $this.query[id.toString()] = function(resp) {
      delete $this.query[id.toString()];
      a(resp);
    }

    $this.conn.post({
      type: "query",
      id:   id.toString()
    })
  })
}

module.exports = Webdaemon;