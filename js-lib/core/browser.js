function WebdaemonBrowser(webd, onmessage) {
  this.sw = new SharedWorker(file + "-sharedworker.js");
  this.sw.port.addEventListener("message", onmessage);

  window.addEventListener("beforeunload", function() {
    this.sw.port.postMessage({
      type: "decommission"
    });
  });

  sw.port.start();
}

WebdaemonBrowser.prototype.post = function(msg) {
  this.sw.port.postMessage(msg);
}

module.exports = WebdaemonBrowser;
