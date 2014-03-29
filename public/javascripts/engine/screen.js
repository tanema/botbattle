function Screen(canvas, tilewidth, tileheight){
  this.canvas = canvas;
  this.tilewidth = tilewidth;
  this.tileheight = tileheight;
}

Screen.prototype.update = function(){
  this.maintain_aspect_ratio()
  this.width  = this.canvas.width;
  this.height = this.canvas.height;
  this.tilesX = this.canvas.width  / this.tilewidth;
  this.tilesY = this.canvas.height / this.tileheight;
}

//keep the width constant and change the height to reflect changes to the screen
Screen.prototype.maintain_aspect_ratio = function(){
  var w = window.innerWidth
    || document.documentElement.clientWidth
    || document.body.clientWidth;

  var h = window.innerHeight
    || document.documentElement.clientHeight
    || document.body.clientHeight;

  var set_height = w/(this.canvas.width/this.canvas.height); 
  if(set_height < h){
    this.canvas.style.width = "100%";
    this.canvas.style.height = set_height+"px";
    this.canvas.style["margin"] = "0px";
  } else {
    var set_width = h/(this.canvas.height/this.canvas.width),
        margin = (w - set_width) / 2;
    this.canvas.style.width = set_width+"px";
    this.canvas.style.height = "100%";
    this.canvas.style["margin"] = "0px "+margin+"px 0px "+margin+"px";
  }
}
