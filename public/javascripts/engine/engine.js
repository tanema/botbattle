window.requestAnimationFrame = window.requestAnimationFrame       ||
                               window.webkitRequestAnimationFrame ||
                               window.mozRequestAnimationFrame    ||
                               function(callback){window.setTimeout(callback, 1000 / 60);}; 

function Grouter(canvas_el, map_src){
  if(!this.canvasIsSupported()){
    alert("Your browser does not support this game.");
    return;
  }

  var _this = this;
  this.canvas = document.getElementById(canvas_el);
  this.ctx = this.canvas.getContext('2d');
  this.ctx.canvas = this.canvas;

  this.canvas.style.position = "absolute";
  this.canvas.style.top = 0;
  this.canvas.style.left =0;

  this.fps_el = document.getElementById("fps"); 
  if(this.fps_el){
    this.fps = 0;
    this.fps_timer = setInterval(function(){_this.updateFPS()}, 1000);
  }

  window.addEventListener('load', function(){
    _this.load_map(map_src);
  });
}

Grouter.gameloop = function(method, scope){
  var startTime = window.mozAnimationStartTime || Date.now(),
      stop = false;
      looper = {
        stop: function(){
          stop = true;
        },
        start: function(){
          stop = false;
          (function _loop(timestamp){
            //calculate difference since last repaint
            var drawStart = (timestamp || Date.now()),
                deltatime = drawStart - startTime;
            method.call(scope, deltatime)
            startTime = drawStart;
            if(!stop){
              setTimeout(function() {requestAnimationFrame(_loop) }, 1000/30);
            }
          })()
        }
      }
  looper.start();
  return looper;
}

Grouter.resolutions = [
  [640, 480],
  [800, 600],
  [1024, 768],
  [1280, 720],
  [1360, 768],
  [1366, 768],
]
Grouter.prototype.set_resolution = function(i){
  this.canvas.width = Grouter.resolutions[i][0];
  this.canvas.height = Grouter.resolutions[i][1];
}

Grouter.prototype.load_map = function(map_src){
  var _this = this;
  if(this.map){
    this.unload_map();
  }
  this.socket = new SocketIO("/ws");
  this.loaded = false;
  this.map = new Map(map_src, this);
  this.map.load(function(map){
    var tile_width = map.spritesheet.tile_width,
        tile_height = map.spritesheet.tile_height,
        screen = _this.ctx.screen = new Screen(_this.canvas, tile_width, tile_height);
    map.camera = _this.ctx.camera = new Camera(screen, tile_width, tile_height, map.properties.tiles_overflow);
    _this.loaded = true;
    _this.gameloop = Grouter.gameloop(_this.draw, _this);
  });
};

Grouter.prototype.draw = function(deltatime){
  if(!this.loaded){return;}
  //clear last frame
  this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

  //update screen
  this.ctx.screen.update()

  // draw down the hierarchy starting at the map
  this.map.draw(this.ctx, deltatime);

  //increments frame for fps display
  if(this.fps_timer){
    this.fps++;
  }
};

Grouter.prototype.updateFPS = function(){
  this.fps_el.innerHTML = this.fps || 0;
  this.fps = 0;
},

Grouter.prototype.canvasIsSupported = function (){
  var elem = document.createElement('canvas');
  return !!(elem && elem.getContext && elem.getContext('2d'));
};

Grouter.prototype.getSocketId = function () {
  return this.socket.socket.sessionid;
}

Grouter.normalize_coord = function(h, j){
  return Math.floor(((2*j)+(h%j))%j)
}

Grouter.merge_objects = function(obj1, obj2) {
  for (var p in obj2) {
    if (obj1[p] && obj2[p].constructor == Object) {
      obj1[p] = Grouter.merge_objects(obj1[p], obj2[p]);
    } else {
      obj1[p] = obj2[p];
    }
  }
  return obj1;
}

Grouter.getJSON = function(url, cb) {
  var xmlhttp = new XMLHttpRequest();
  xmlhttp.onreadystatechange = function() {
    if (xmlhttp.readyState==4 && xmlhttp.status==200) {
      cb(JSON.parse(xmlhttp.responseText))
    } else if (xmlhttp.readyState==4){
      cb()
    }
  }
  xmlhttp.open("GET", url, true);
  xmlhttp.send();
}

Grouter.bind_event = function(event_names, cb, scope){
  var events = event_names.split(" ")
  for(var i = 0; i < events.length; i++){
    document.addEventListener(events[i], function(e){
      try{
        cb.call(scope, e)
      }catch(e){
        this.removeEventListener(events[i], arguments.callee);
      }
    });
  }
}

Grouter.fire_event = function(name, extra) {
  document.dispatchEvent(new CustomEvent(name, extra))
}
