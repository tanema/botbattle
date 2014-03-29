function Map(map_src, engine){
  this.map_src = map_src;
  this.engine = engine;
  this.layers = {};
  this.camera = null;
  this.sprites = {};
  this.audio_manager = new AudioManager();
  this.name = map_src.substring(map_src.lastIndexOf("/")+1, map_src.lastIndexOf("."));
}

Map.prototype.load = function (next){
  var _this = this, i;

  console.log("["+ _this.map_src + "] Loading");
  Grouter.getJSON(this.map_src, function(map_data){
    _this.properties = map_data.properties || {};
    _this.orientation = map_data.orientation;
    _this.tilewidth = map_data.tilewidth;
    _this.tileheight = map_data.tileheight;
    _this.width = map_data.width;
    _this.height = map_data.height;

    if(_this.properties.music){
      var sound = _this.audio_manager.load_music(_this.properties.music);
      _this.audio_manager.loop(sound);
    }

    console.log(" → loading " + map_data.tilesets.length + " tileset(s)");
    _this.spritesheet = new SpriteSheet(map_data.tilewidth, map_data.tileheight);
    _this._load_tileset(map_data.tilesets, function(){
      console.log(" → spritesheet loaded: " + _this.spritesheet.loaded());
      console.log(" → setting up " + map_data.layers.length + " layer(s)");
      _this._load_layer(map_data.layers, function() {   
        _this.register_socket_events();
        console.log(" → finished loading map data");
        if(next){// everything loaded so continue 
        next(_this);
      }
    });
  });
  });
};

Map.prototype.register_socket_events = function(){
  var self = this;
  console.log(" → connecting to sockets");
  this.socket = this.engine.socket;
  this.socket.on("register",      function(bot_id, x, y, rot, name){ self.on_register(bot_id, x, y, rot, name)})
  this.socket.on("kill",          function(bot_id){ self.on_kill(bot_id)})
  this.socket.on("rotate left",   function(bot_id){ self.on_rotate_left(bot_id)})
  this.socket.on("rotate right",  function(bot_id){ self.on_rotate_right(bot_id)})
  this.socket.on("move forward",  function(bot_id){ self.on_move_forward(bot_id)})
  this.socket.on("move backward", function(bot_id){ self.on_move_backward(bot_id)})
  this.socket.on("fire gun",      function(bot_id){ self.on_fire_gun(bot_id)})
  this.socket.on("fire cannon",   function(bot_id){ self.on_fire_cannon(bot_id)})
  this.socket.on("scan",          function(bot_id){ self.on_scan(bot_id)})
};

Map.prototype.loaded = function (){
  return this.spritesheet.loaded();
};

Map.prototype._load_tileset = function(tilesets, next){
  if(tilesets.length === 0){return next();}
  var _this = this;
  this.spritesheet.add_image(tilesets[0], function(){
    tilesets.shift();
    _this._load_tileset(tilesets, next);
  });
};

Map.prototype._load_layer = function(layers, next){
  if(layers.length === 0){return next();}
  var _this = this;
  new Layer(layers[0], _this, function(layer){
    _this.layers[layer.name] = layer;
    layers.shift();
    _this._load_layer(layers, next);
  });
};

Map.prototype.at = function(x, y, group){
  var results = {tiles: [], sprites: [], actors: []};

  var layer_name;
  for(layer_name in this.layers){
    var layer = this.layers[layer_name];

    if(!layer.visible || group != layer.group){
      continue;
    }

    if(layer.is_tilelayer()){
      var tile = this.spritesheet.get(layer.get_tile_index(x, y));
      if(tile){
        results.tiles.push(tile);
      }
    }else if(layer.is_objectgroup()){
      var object_name, object;
      for(object_name in layer.sprites){
        object = layer.sprites[object_name];
        if(object && object.x === x && object.y === y){
          results.sprites.push(object);
        }
      }
      for(object_name in layer.actors){
        object = layer.actors[object_name];
        if(object && object.x === x && object.y === y){
          results.actors.push(object);
        }
      }
    }
  }

  return results;
};

Map.prototype.draw = function (ctx, deltatime){
  //default background, using css
  ctx.canvas.style.background = this.properties.background

  var layer_name;
  for(layer_name in this.layers){
    this.layers[layer_name].draw(ctx, deltatime);
  }
};

Map.prototype.on_register = function(id, x, y, rot, name){
  console.log(" → spawning " + id + ":" + name + " at " + x + "," + y + " rotation: " + rot);
  var sprite = new Sprite({
    id: id,
    x: x,
    y: y,
    name: name,
    rotation: rot
  }, this, this.layers["players"])
}
Map.prototype.on_kill = function(bot_id){
  try{
    this.sprites[bot_id].kill()
  } catch(e){}
}

Map.prototype.on_rotate_left = function(bot_id){
  try{
    var bot = this.sprites[bot_id];
    if(bot.rotation == 0){
      bot.rotation = 270
    } else {
      bot.rotation -= 90
    }
  } catch(e){}
}

Map.prototype.on_rotate_right = function(bot_id){
  try{
    var bot = this.sprites[bot_id];
    if(bot.rotation == 270){
      bot.rotation = 0
    } else {
      bot.rotation += 90
    }
  } catch(e){}
}

Map.prototype.on_move_forward = function(bot_id){
  try{
    this.sprites[bot_id].moveForward()
  } catch(e){}
}

Map.prototype.on_move_backward = function(bot_id){
  try{
    this.sprites[bot_id].moveBackward()
  } catch(e){}
}

Map.prototype.on_fire_gun = function(bot_id){
  try{
    this.sprites[bot_id].fireGun()
  } catch(e){}
}

Map.prototype.on_fire_cannon = function(bot_id){
  try{
    this.sprites[bot_id].fireCannon()
  } catch(e){}
}

Map.prototype.on_scan = function(bot_id){
  try{
    this.sprites[bot_id].scan()
  } catch(e){}
}
