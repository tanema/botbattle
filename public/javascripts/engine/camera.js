function Camera(screen, tile_width, tile_height, tiles_overflow){
  this.screen = screen;
  this.x = 0;
  this.y = 0;
  this.tile_width = tile_width;
  this.tile_height = tile_height;
  this.overflowTile = tiles_overflow || 5;
}

Camera.prototype.set = function(x, y) {
  var pos = this.get_center_on(x,y);
  this.x = pos.x;
  this.y = pos.y;
};

Camera.prototype.get_center_on = function(x,y){
  return {
    x: x - (this.screen.width  - this.tile_width) / (this.tile_width * 2),
    y: y - (this.screen.height - this.tile_height) / (this.tile_height * 2)
  }
}

Camera.prototype.pan_to = function(x, y, cb){
  Animation.ease(this, this.get_center_on(x, y), {callback: cb})
}

Camera.prototype.left = function(){
  return this.x - this.overflowTile;
};

Camera.prototype.right = function(){
  return this.x + this.screen.tilesX + this.overflowTile;
};

Camera.prototype.top = function(){
  return this.y - this.overflowTile;
};

Camera.prototype.bottom = function(){
  return this.y + this.screen.tilesY + this.overflowTile;
};

Camera.prototype.isInside = function(x, y, layer) {
  if(x > this.left() && x < this.right() && y > this.top() && y < this.bottom()){
    return [x,y]
  } else {
    var new_x, new_y;
    if(x > this.left() && x < this.right()){
      new_x = x;
    } else if((x+layer.width) > this.left() && (x+layer.width) < this.right()){
      new_x = x+layer.width;
    } else if((x-layer.width) > this.left() && (x-layer.width) < this.right()){
      new_x = x-layer.width;
    }
    if(y > this.top() && y < this.bottom()){
      new_y = y;
    } else if((y+layer.height) > this.top() && (y+layer.height) < this.bottom()){
      new_y = y+layer.height;
    } else if((y-layer.height) > this.top() && (y-layer.height) < this.bottom()){
      new_y = y-layer.height;
    }
    if(new_x && new_y){
      return [new_x, new_y]
    }
  }
  return false
};
