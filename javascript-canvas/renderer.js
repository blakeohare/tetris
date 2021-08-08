const Renderer = (() => {

  const keyCodeLookup = {
    ArrowLeft: "left",
    ArrowRight: "right",
    ArrowUp: "up",
    ArrowDown: "down",
    Space: "space",
  };

  let ctx = null;
  let ctxWidth = null;
  let ctxHeight = null;

  let hexLookup = [];
  let hexLetters = '0123456789abcdef';
  for (i = 0; i < 256; ++i) {
    hexLookup.push(hexLetters.charAt(i >> 4) + hexLetters.charAt(i & 15));
  }
  let hexLookupHash = hexLookup.map(v => '#' + v);
  let toHex = (r, g, b) => {
    return hexLookupHash[r] + hexLookup[g] + hexLookup[b];
  };

  let gfx = {
    fill: (r, g, b) => {
      ctx.fillStyle = toHex(r, g, b);
      ctx.fillRect(0, 0, ctxWidth, ctxHeight);
    },
    rectangle: (x, y, width, height, r, g, b) => {
      ctx.fillStyle = toHex(r, g, b);
      ctx.fillRect(x, y, width, height);
    },
  };

  return {
    start: (width, height, update, render, canvasHost) => {
      ctxWidth = width;
      ctxHeight = height;

      let canvas = document.createElement('canvas');
      canvas.width = width;
      canvas.height = height;
      canvas.style.width = width + 'px';
      canvas.style.height = height + 'px';
      canvasHost.append(canvas)

      let virtualCanvas = document.createElement('canvas');
      virtualCanvas.width = width;
      virtualCanvas.height = height;
      ctx = virtualCanvas.getContext('2d');

      let events = [];

      let pressedKeys = {};
      let handleKey = (keyCode, down) => {
        let keyName = keyCodeLookup[keyCode] || null;
        if (!down && !pressedKeys[keyName]) return;

        events.push({
          key: keyName,
          down
        });
        pressedKeys[keyName] = down;
      };

      window.addEventListener('keydown', e => handleKey(e.code, true));
      window.addEventListener('keyup', e => handleKey(e.code, false));

      let renderCounter = 0;

      window.setInterval(() => {
        let eventsThisFrame = events.slice(0);
        while (events.length) events.pop();
        update(eventsThisFrame);
        render(gfx, renderCounter++);

        canvas.getContext('2d').drawImage(virtualCanvas, 0, 0, width, height);

      }, 1000 / 60);
    },
  }
})();
