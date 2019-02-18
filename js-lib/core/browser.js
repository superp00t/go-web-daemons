function WebdaemonBrowser(file, onmessage) {
  this.sw = new SharedWorker(file + "-sharedworker.js");
  this.sw.port.addEventListener("message", function(evt) {
    onmessage(JSON.parse(evt.data));
  });

  window.addEventListener("beforeunload", function() {
    this.sw.port.postMessage({  
      type: "decommission"
    });
  });

  this.sw.port.start();
}

WebdaemonBrowser.prototype.post = function(msg) {
  this.sw.port.postMessage(JSON.stringify(msg));
}

module.exports = WebdaemonBrowser;
