window.onload = function() {
  var canvas,
    ctx,
    nodes,
    links,
    render,
    getMousePosFromEvent,
    getNodeByPos,
    dragNode,
    dragPoint;

  canvas = document.getElementById('viewport');
  ctx = canvas.getContext('2d');

  // Массив узлов, каждый узел задается объектом с координатами и id картинки для отображения.
  nodes = [{
    x: 300,
    y: 50,
    id: 'picture_a'
  }, {
    x: 200,
    y: 100,
    id: 'picture_b'
  }, {
    x: 400,
    y: 100,
    id: 'picture_c'
  }];

  // Массив связей между узлами. Каждый элемент это массив из двух элементов - индексов узлов.
  links = [
    [0, 1],
    [0, 2]
  ];

  // Отрисовка канвы.
  render = function() {

    // Очищаем канву.
    ctx.fillStyle = "#f3f3f3";
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    // Рисуем связи между узлами (раньше чем сами узлы, чтобы они отображались позади узлов).
    links.forEach(function(link) {
      var i0 = link[0],
        i1 = link[1];
      ctx.fillStyle = "#000";
      ctx.lineWidth = 2;
      ctx.beginPath();
      ctx.moveTo(nodes[i0].x, nodes[i0].y);
      ctx.lineTo(nodes[i1].x, nodes[i1].y);
      ctx.stroke();
    });

    // Рисуем узлы.
    nodes.forEach(function(node) {
      var img = document.getElementById(node.id),
        halfWidth = img.naturalWidth / 2,
        halfHeight = img.naturalHeight / 2;
      ctx.drawImage(img, node.x - halfWidth, node.y - halfHeight);
    });
  };

  // Получает из события мыши координаты, относительно левого верхнего угла канвы.
  getMousePosFromEvent = function(evt) {
    var rect = canvas.getBoundingClientRect();
    return {
      x: evt.clientX - rect.left,
      y: evt.clientY - rect.top
    };
  };

  // Находит узел, находящийся по заданой координате на канве.
  getNodeByPos = function(pos) {
    var result;
    nodes.forEach(function(node) {
      var img = document.getElementById(node.id),
        halfWidth = img.naturalWidth / 2,
        halfHeight = img.naturalHeight / 2;
      if ((pos.x >= node.x - halfWidth) && (pos.x < node.x + halfWidth) && (pos.y >= node.y - halfHeight) && (pos.y < node.y + halfHeight)) {
        result = node;
      }
    });
    return result;
  };

  // При нажатии кнопки мыши находим узел по которому было нажатие,
  // запоминаем его в dragNode для дальнейшего использования,
  // в dragPoint запоминаем по какому месту узла была нажата кнопка мыши.
  canvas.addEventListener('mousedown', function(event) {
    var pos = getMousePosFromEvent(event);
    dragNode = getNodeByPos(pos);
    if (dragNode !== undefined) {
      dragPoint = {
        x: pos.x - dragNode.x,
        y: pos.y - dragNode.y
      }
    }
  }, false);

  // При отпускании кпнопки мыши забываем текущий перетаскиваемый узел.
  canvas.addEventListener('mouseup', function() {
    dragNode = undefined;
  }, false);

  // При движении мыши, если есть перетаскиваемый узел, двигаем его и перерисовываем канву.
  canvas.addEventListener('mousemove', function(event) {
    var pos;
    if (dragNode !== undefined) {
      pos = getMousePosFromEvent(event);
      dragNode.x = pos.x - dragPoint.x;
      dragNode.y = pos.y - dragPoint.y;
      render();
    }
  }, false);

  render();
};
