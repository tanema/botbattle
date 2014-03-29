function SocketIO(subpath){
  this.path = "ws://"+window.location.host;
  if(subpath){
    this.path += subpath
  }
  this.$events = {};
  this.$callbacks = {};
  this.connect();
}

SocketIO.prototype.connect = function() {
  var self = this;
  this.socket = new WebSocket(this.path);  
  this.socket.onopen =    function(e){ self.onOpen(e) };
  this.socket.onmessage = function(e){ self.onMessage(e) };
  this.socket.onclose =   function(e){ self.onClose(e) };
  this.socket.onerror =   function(e){ self.onError(e) };
}

SocketIO.prototype.onOpen = function(e) {
}

SocketIO.prototype.onClose = function(e) {
  var self = this;
  setTimeout(function(){
    self.connect();
  }, 1000)
};

SocketIO.prototype.onError = function(e) {
}

SocketIO.prototype.onMessage = function(resp) {
  var message = JSON.parse(resp.data),
      handlers = this.$events[message.event_name],
      callbacks = this.$callbacks[message.event_name];

  if (handlers) {
    for (var i = 0, l = handlers.length; i < l; i++) {
      handlers[i].apply(this, message.event_data);
    }
  }
  if (callbacks) {
    for (var i = 0, l = callbacks.length; i < l; i++) {
      callbacks[i].apply(this, message.event_data);
    }
    delete this.$callbacks[message.event_name];
  }

}

SocketIO.prototype.emit = function() {
  if(this.socket.readyState == WebSocket.CLOSED){
    console.log("the websocket has lost connection")
    return
  }

  var args = Array.prototype.slice.call(arguments),
      event_name = args.shift();

  if(typeof args[args.length - 1] === "function"){
    var cb = args.pop();
    if (!this.$callbacks[event_name + " response"]) {
      this.$callbacks[event_name + " response"] = [cb];
    } else {
      this.$callbacks[event_name + " response"].push(cb);
    } 
  }

  var message = {
    event_name: event_name,
    event_data: args
  };

  this.socket.send(JSON.stringify(message));
}

SocketIO.prototype.on = function(event_name, fn) {
  if (!this.$events[event_name]) {
    this.$events[event_name] = [fn];
  } else {
    this.$events[event_name].push(fn);
  } 
}
