function AudioManager(type){
  this.step_size = 0.05;

  this.music_volume = localStorage["music_vol"] || 0.5;
  this.sfx_volume   = localStorage["sfx_vol"] || 0.5;

  this.music = {};
  this.sfx = {};

  this._bind_change_events();
}

AudioManager.prototype.load_music = function(src){
  var sound = new Audio(src),
      filename = sound.src.substring(sound.src.lastIndexOf('/')+1);
  filename = filename.substring(0, filename.lastIndexOf('.'));

  sound.volume = this.music_volume;
  this.music[filename] = sound;

  return filename;
};

AudioManager.prototype.load_sfx = function(src){
  var sound = new Audio(src),
      filename = sound.src.substring(sound.src.lastIndexOf('/')+1);
  filename = filename.substring(0, filename.lastIndexOf('.'));

  sound.volume = this.sfx_volume;
  this.sfx[filename] = sound;

  return filename;
};

AudioManager.prototype.play = function(sound_name){
  this.getSound(sound_name).play();
};

AudioManager.prototype.getSound = function(sound_name){
  return this.music[sound_name] || this.sfx[sound_name];
};

AudioManager.prototype.stop = function(sound_name){
  this.pause(sound_name);
  this.currentTime = 0;
};

AudioManager.prototype.pause = function(sound_name){
  this.getSound(sound_name).pause();
};

AudioManager.prototype.loop = function(sound_name){
  var sound = this.getSound(sound_name);
  if (typeof sound.loop === 'boolean'){
    sound.loop = true;
  }else{
    sound.addEventListener('ended', function() {
      this.currentTime = 0;
      this.play();
    }, false);
  }
  sound.play();
};

AudioManager.prototype.change_volume = function(vol, type){
  localStorage["prev_"+type+"_vol"] = localStorage[type+"_vol"] || 0.5;
  this[type+"_volume"] = localStorage[type+"_vol"] = vol;
  for(var sound_name in this[type]){
    this[type][sound_name].volume = this[type+"_volume"];
  }
};

AudioManager.prototype._bind_change_events = function(){
  var _this = this;
  //music events
  document.addEventListener("music_off", function(e){
    _this.change_volume(0, "music");
  })
  document.addEventListener("music_on", function(){
    _this.change_volume(localStorage["prev_music_vol"] && localStorage["prev_music_vol"] != "0" ? localStorage["prev_music_vol"] : 0.5, "music");
  })
  document.addEventListener("music_vol_down", function(){
    _this.change_volume(_this.music_volume - _this.step_size, "music");
  })
  document.addEventListener("music_vol_up", function(){
    _this.change_volume(_this.music_volume + _this.step_size, "music");
  })
  document.addEventListener("music_vol_change", function(e, vol){
    _this.change_volume(vol/100, "music");
  });

  //sfx events
  document.addEventListener("sfx_off", function(){
    _this.change_volume(0, "sfx");
  })
  document.addEventListener("sfx_on", function(){
    _this.change_volume(localStorage["prev_sfx_vol"] && localStorage["prev_sfx_vol"] != "0" ? localStorage["prev_sfx_vol"] : 0.5, "sfx");
  })
  document.addEventListener("sfx_vol_down", function(){
    _this.change_volume(_this.sfx_volume - _this.step_size, "sfx");
  })
  document.addEventListener("sfx_vol_up", function(){
    _this.change_volume(_this.sfx_volume + _this.step_size, "sfx");
  })
  document.addEventListener("sfx_vol_change", function(e, vol){
    _this.change_volume(vol/100, "sfx");
  });
};
