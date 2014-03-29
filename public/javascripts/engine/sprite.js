function Sprite(display_object_options, map, layer){
  display_object_options = display_object_options || {};

  this.map = map;
  this.id = display_object_options.id;
  this.name = display_object_options.name;
  this.rotation = display_object_options.rotation || 0;
  this.layer = layer;
  this.x = display_object_options.x || 0;
  this.y = display_object_options.y || 0;
  this.health = 100;

  this.is_scanning = false;
  this.is_firing_cannon = false;
  this.is_firing_gun = false;

  this.initalize_properties();
  this.animationloop = Grouter.gameloop(this.animate, this);

  map.sprites[this.id] = this;
  layer.sprites[this.id] = this;
}

Sprite.prototype.initalize_properties = function(){
  this.width = this.map.spritesheet.tile_width;
  this.height = this.map.spritesheet.tile_height;

  this.currentMovement = "down";
  this.frame_time = 0;
  this.movementIndex = 0;

  this.sprite_index = 47;
  this.speed =  200;
  this.spritesheet = this.map.spritesheet;
};

Sprite.prototype.draw = function(ctx){
  var draw_x = (this.x * this.width),
      draw_y = (this.y * this.height);

  var sprite = this.spritesheet.get(this.sprite_index),
      draw_frame = sprite ? sprite.img : null;

  if(draw_frame){
    ctx.save(); 
    ctx.translate(draw_x, draw_y); 
    var middle_x = (this.width/2),
        middle_y = (this.height/2);
    ctx.translate(middle_x, middle_y); 
    ctx.rotate(Math.PI / 180 * (this.rotation - 90)); 
    ctx.drawImage(draw_frame, -middle_x, -middle_y)
    ctx.restore();
  }
  if(this.is_scanning){
    this.drawScan(ctx)
  }
  if(this.is_firing_cannon){
    this.drawCannon(ctx)
  }
  if(this.is_firing_gun){
    this.drawGun(ctx)
  }
};

Sprite.prototype.drawScan = function(ctx){
  this._drawWeapon(function(start_x, start_y, end_x, end_y){
    ctx.save(); 
    ctx.beginPath();
    ctx.moveTo(start_x, start_y);
    ctx.lineTo(end_x, end_y);

    ctx.lineWidth = 10;
    ctx.strokeStyle = 'rgba(0,255,0,0.5);';

    ctx.stroke();
    ctx.restore();
  })
}

Sprite.prototype.drawCannon = function(ctx){
  this._drawWeapon(function(start_x, start_y, end_x, end_y){
    ctx.save(); 
    ctx.beginPath();
    ctx.moveTo(start_x, start_y);
    ctx.lineTo(end_x, end_y);

    ctx.lineWidth = 10;
    ctx.strokeStyle = 'rgba(255,0,0,0.7)';

    ctx.stroke();
    ctx.restore();
  }, true)
}
Sprite.prototype.drawGun = function(ctx){
  this._drawWeapon(function(start_x, start_y, end_x, end_y){
    ctx.save(); 
    ctx.beginPath();
    ctx.setLineDash([5])
    ctx.moveTo(start_x, start_y);
    ctx.lineTo(end_x, end_y);
    
    ctx.lineWidth = 5;
    ctx.strokeStyle = 'rgba(255,0,0,0.7)';

    ctx.stroke();
    ctx.restore();
  }, true)
}

Sprite.prototype._drawWeapon = function(draw_cb, first_item){
  var x = (this.x * this.width),
      y = (this.y * this.height),
      item = first_item ? this.lookingAt()[0] : null;

  switch(this.rotation){
  case 90:
    draw_cb(
      x + (this.width/2), y, 
      x + (this.width/2), (item ? (item.y + this.height) : 0)
    )
    break;
  case 270:
    draw_cb(
      x + (this.width/2), y+this.height,
      x + (this.width/2), (item ? item.y : 12*this.width)
    )
    break;
  case 0:
    draw_cb(
      x, y+(this.height/2), 
      (item ? (item.x + this.width) : 0), y+(this.height/2)
    )
    break;
  case 180:
    draw_cb(
      x+this.width, y+(this.height/2), 
      (item ? item.x : 24*this.height), y+(this.height/2)
    )
    break;
  }
}

Sprite.prototype.set = function(x, y){
  this.is_moving = false
  this.x = x;
  this.y = y;
};

Sprite.prototype.kill = function(){
  delete this.layer.sprites[this.id];
  delete this.map.sprites[this.id];
};

//@OVERRIDE this just make sure the displayable is facing the speaker/actor
Sprite.prototype.unload = function(){
  this.animationloop.stop();
  //call Super
  this.constructor.prototype.unload.call(this);
};

Sprite.prototype.moveForward = function(){
  switch(this.rotation){
  case 90:
    if(this.y > 0) {
      this.y--
    }
    break;
  case 270:
    if(this.y < 12) {
      this.y++
    }
    break;
  case 0:
    if(this.x > 0) {
      this.x--
    }
    break;
  case 180:
    if(this.x < 24) {
      this.x++
    }
    break;
  }
};

Sprite.prototype.moveBackward = function(){
  switch(this.rotation){
  case 90:
    if(this.y < 12){
      this.y++
    }
    break;
  case 270:
    if(this.y > 0) {
      this.y--
    }
    break;
  case 0:
    if(this.x < 24) {
      this.x++
    }
    break;
  case 180:
    if(this.x > 0) {
      this.x--
    }
    break;
  }
};

Sprite.prototype.scan = function(){
  this.is_scanning = true;
  var self = this;
  setTimeout(function(){
    self.is_scanning = false;
  }, 500)
};

Sprite.prototype.fireCannon = function(){
  this.is_firing_cannon = true;
  var self = this;
  setTimeout(function(){
    self.is_firing_cannon = false;
  }, 3000)
};

Sprite.prototype.fireGun = function(){
  this.is_firing_gun = true;
  var self = this;
  setTimeout(function(){
    self.is_firing_gun = false;
  }, 1000)
};

Sprite.prototype.lookingAt = function(){
  var result = [],
      bots = this.map.sprites;

  for(var i = 0; i < bots.length; i++){
    switch(this.rotation){
    case 90:
      if(this.y > bot.y && this.x == bot.x) {
        result.push(bot)
      }
    case 270:
      if(this.y < bot.y && this.x == bot.x) {
        result.push(bot)
      }
    case 0:
      if(this.y == bot.y && this.x > bot.x) {
        result.push(bot)
      }
    case 180:
      if(this.y == bot.y && this.x < bot.x) {
        result.push(bot)
      }
    }
  }

  var self = this;
  result.sort(function(b1, b2){
    switch(self.rotation){
    case 90:
      return b1.y > b2.y
    case 270:
      return b1.y < b2.y
    case 0:
      return b1.x > b2.x
    case 180:
      return b1.x < b2.x
    }
  })
  return result
}

Sprite.prototype.animate = function(deltatime){
};
