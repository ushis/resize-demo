(function() {

  //
  //
  //
  var Util = {

    //
    el: function(name, attributes) {
      var el = document.createElement(name);

      for (attr in attributes) {
        el.setAttribute(attr, attributes[attr]);
      }
      return el;
    },

    //
    img: function(src, callback) {
      var img = new Image();
      img.addEventListener('load', callback, false);
      img.src = src;
      return img;
    },

    //
    stop: function(e) {
      e.stopPropagation();
      e.preventDefault();
    },

    //
    basename: function(path) {
      return path.split('/').pop();
    },

    //
    removeExt: function(path) {
      return path.split('.').shift();
    }
  };

  //
  //
  //
  var Original = function(el, results) {
    this.el = el;
    this.results = results;

    this.el.addEventListener('dragenter', Util.stop);
    this.el.addEventListener('dragleave', Util.stop);
    this.el.addEventListener('dragover', Util.stop);

    var that = this;

    this.el.addEventListener('drop', function(e) {
      Util.stop(e);

      if (e.dataTransfer.files.length > 0) {
        that.send(e.dataTransfer.files.item(0));
      }
    });

    this.refresh();
  };

  //
  Original.prototype.send = function(file) {
    if ( ! /^image\/(?:png|jpe?g|gif)$/i.test(file.type)) {
      return alert('This is not an image file!');
    }
    var that = this;
    var gutter = Util.el('div', {class: 'gutter'});
    that.el.appendChild(gutter);

    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/', true);
    xhr.reponseType = 'text';

    xhr.upload.addEventListener('progress', function(e) {
      gutter.style.top = 100 - ((e.total > 0) ? (e.loaded * 100) / e.total : 0) + '%';
    });

    xhr.addEventListener('load', function(_) {
      if (this.status != 200) {
        return alert(this.responseText);
      }
      that.display(this.responseText);
    });

    xhr.send(file)
  }

  //
  Original.prototype.display = function(src) {
    var el = this.el;
    el.innerHTML = '';
    el.classList.add('loading');

    el.appendChild(Util.img(src, function() {
      el.classList.remove('loading');
    }));

    this.results.update(src);
    window.history.pushState(null, 'Resize Demo', Util.removeExt(Util.basename(src)));
  };

  //
  Original.prototype.refresh = function() {
    var path = window.location.pathname.slice(1);

    if (path.length > 0) {
      this.display(['/images', 'orig', path[0], path[1], path + '.png'].join('/'));
    }
  }

  //
  //
  //
  var Results = function(el) {
    this.el = el;
    this.src = '';
    this.images = el.querySelectorAll('#images li div.img');
    this.activeTab = el.querySelector('#tabs li.active');

    var that = this;
    var tabs = that.el.querySelectorAll('#tabs li');

    for (var i = 0, len = tabs.length; i < len; i++) {
      tabs[i].addEventListener('click', function() {
        that.activeTab.classList.remove('active');
        that.activeTab = this;
        this.classList.add('active');
        that.display();
      });
    }
  };

  //
  Results.prototype.update = function(url) {
    this.src = url.split('/').slice(-3).join('/');
    this.el.classList.add('active');
    this.display();
  };

  //
  Results.prototype.display = function() {
    var width = this.activeTab.dataset.width;

    for (var i = 0, len = this.images.length; i < len; i++) {
      var img = this.images[i];
      var src = ['images', width, img.dataset.interpolation, this.src];
      img.innerHTML = '';

      var li = img.parentNode;
      li.style.width = width + 'px';
      li.classList.add('loading');

      img.appendChild(Util.img(src.join('/'), function() {
        this.parentNode.parentNode.classList.remove('loading');
      }));
    }
  };

  // Let's go.
  var results = new Results(document.querySelector('#results'));
  var orig = new Original(document.querySelector('#orig'), results);
  window.addEventListener('popstate', function() { orig.refresh(); });
}).call(this);
