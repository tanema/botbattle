function Explosion(x, y, sprite){
  this.x = x;
  this.y = y;
  this.sprite = sprite;
  
  this.sprite_end = 45;

  this.explosions = [];
  function rand(min, max) {
    return parseInt(Math.random() * (max - min) + min);
  }

  for(var i = 0; i < 20; i++){
    this.explosions.push({
      index: rand(0, 5),
      x: x + rand(-10, 10),
      y: y + rand(-10, 10),
      rot: rand(0, 270),
    })
  }
  this.frame_time = 0;
  this.animation_speed = 45
  this.animationloop = Grouter.gameloop(this.animate, this);
}

Explosion.prototype.animate = function(deltatime){
  if(this.frame_time < 0){
    this.frame_time = 0
  }
  if((this.frame_time += deltatime) >= this.animation_speed){
    for(var i=0; i < this.explosions.length; i++){
      this.explosions[i].index++
      this.explosions[i].rot += 5
    }
    this.frame_time = 0;
  }
}

Explosion.prototype.draw = function(ctx){
  for(var i=0; i < this.explosions.length; i++){
    var exp = this.explosions[i],
        sprite = this.sprite.spritesheet.get(exp.index),
        draw_frame = sprite ? sprite.img : null;

    if(draw_frame && exp.index <= this.sprite_end){
      ctx.save(); 
      ctx.translate(exp.x, exp.y); 
      var middle_x = (this.sprite.width/2),
          middle_y = (this.sprite.height/2);
      ctx.translate(middle_x, middle_y); 
      ctx.rotate(Math.PI / 180 * exp.rot); 
      ctx.drawImage(draw_frame, -middle_x, -middle_y)
      ctx.restore();
    }
  }
}
