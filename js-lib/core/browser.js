
function WebdaemonBrowser(webd, onmessage) {
  var $this = this;
  this.sw = new SharedWorker(file + "-sharedworker.js");
  sw.port.addEventListener("message", onmessage);

  window.addEventListener("beforeunload", function() {
    sw.port.postMessage({
      type: "decommission"
    });
  });

  sw.port.start();
}

WebdaemonBrowser.prototype.post = function(msg) {
  this.sw.postMessage(msg);
}

module.exports = WebdaemonBrowser;
