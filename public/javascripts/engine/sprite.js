function Sprite(display_object_options, map, layer){
  display_object_options = display_object_options || {};

  this.map = map;
  this.id = display_object_options.id;
  this.name = display_object_options.name;
  this.rotation = display_object_options.rotation || 0;
  this.layer = layer;
  this.x = display_object_options.x || 0;
  this.y = display_object_options.y || 0;
  this.health = display_object_options.health || 100;

  this.is_scanning = false;
  this.is_firing_cannon = false;
  this.is_firing_gun = false;

  this.initalize_properties();

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

  if(draw_frame && !this.explosion){
    ctx.save(); 
    ctx.translate(draw_x, draw_y); 
    var middle_x = (this.width/2),
        middle_y = (this.height/2);
    ctx.translate(middle_x, middle_y); 
    ctx.rotate(Math.PI / 180 * (this.rotation - 90)); 
    ctx.drawImage(draw_frame, -middle_x, -middle_y)
    ctx.restore();

    ctx.save()
    ctx.font = '12px pokemon';
    ctx.fillStyle = 'white';
    ctx.fillText(this.name, draw_x+5, draw_y+15);
    ctx.restore();
  } else if(this.explosion){
    this.explosion.draw(ctx)
  }

  if(this.is_scanning && !this.explosion){
    this.drawScan(ctx)
  }
  if(this.is_firing_cannon && !this.explosion){
    this.drawCannon(ctx)
  }
  if(this.is_firing_gun && !this.explosion){
    this.drawGun(ctx)
  }
  if(this.shield && !this.explosion){
    ctx.beginPath();
    ctx.arc(draw_x+(this.width/2), draw_y+(this.height/2), (32/2)+10, 0, 2 * Math.PI, false);
    ctx.fillStyle = 'rgba(0,0,255,0.5)';
    ctx.fill();
  }
};

Sprite.prototype.drawScan = function(ctx){
  this._drawWeapon(function(start_x, start_y, end_x, end_y){
    ctx.save(); 
    ctx.beginPath();
    ctx.moveTo(start_x, start_y);
    ctx.lineTo(end_x, end_y);

    ctx.lineWidth = 10;
    ctx.strokeStyle = 'rgba(0,255,0,0.5)';

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
      x + (this.width/2), (item ? ((item.y * this.height) + this.height/2) : 0)
    )
    break;
  case 270:
    draw_cb(
      x + (this.width/2), y+this.height,
      x + (this.width/2), (item ? ((item.y * this.height) + this.height/2) : 12*this.width)
    )
    break;
  case 0:
    draw_cb(
      x, y+(this.height/2), 
      (item ? ((item.x * this.width) + this.width) : 0), y+(this.height/2)
    )
    break;
  case 180:
    draw_cb(
      x+this.width, y+(this.height/2), 
      (item ? ((item.x * this.width) + this.width) : 24*this.height), y+(this.height/2)
    )
    break;
  }
}

Sprite.prototype.set = function(x, y){
  this.x = x;
  this.y = y;
};

Sprite.prototype.kill = function(){
  var x = ( this.x * this.width),
      y = ( this.y * this.height ),
      self = this;

  this.explosion = new Explosion(x, y, this);
  setTimeout(function(){
    delete self.layer.sprites[self.id];
    delete self.map.sprites[self.id];
  },3000)
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

  for(var bot_id in bots){
    var bot = bots[bot_id];
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
      return b1.y < b2.y
    case 270:
      return b1.y > b2.y
    case 0:
      return b1.x < b2.x
    case 180:
      return b1.x > b2.x
    }
  })
  return result
}
