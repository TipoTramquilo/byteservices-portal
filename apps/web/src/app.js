/* Backend hooks. Set before this script, e.g. window.BYTE_PORTAL={send:"/api/send"}.
   The API identifies the device by the request's IP, so the page sends no device id.
   send   POST {name, id_number:"V18456789", company, phone:"+584125551234"}  -> 2xx {verified:true, redirect?}
   Without an endpoint the portal runs in demo mode.
   Target: ES2020 (docs/explanation/web/es2020.md); `pnpm es-check` validates it. */
(() => {
  const C = window.BYTE_PORTAL || {};
  const DEMO = !C.send;
  const $ = (i) => document.querySelector(`#${i}`);
  const RULE = 40;
  const REDUCED = () => matchMedia('(prefers-reduced-motion: reduce)').matches;
  let data = {};

  /* the byte pattern: drawn once per corner, densest at the bottom-right, azul -> violeta */
  const PAL = [
    '#161D2E',
    '#3A4358',
    '#1A6FB8',
    '#309AE7',
    '#5FD6FF',
    '#309AE7',
    '#5036EA',
    '#5036EA',
  ];
  function pattern(svg) {
    let seed = 11;
    const R = () => {
      seed = (seed * 16807) % 2147483647;
      return seed / 2147483647;
    };
    const CELL = 15;
    const cols = (svg.viewBox.baseVal.width / CELL) | 0;
    const rows = (svg.viewBox.baseVal.height / CELL) | 0;
    let out = '';
    for (let y = 0; y < rows; y++) {
      for (let x = 0; x < cols; x++) {
        const d = (x / cols) * 0.55 + (y / rows) * 0.45;
        if (R() < d * 1.2 - 0.28) {
          const fill = PAL[Math.min(7, (d * 8 * (0.7 + R() * 0.45)) | 0)];
          const o = Math.min(1, 0.3 + d * 0.75).toFixed(2);
          const tw = (R() * 5).toFixed(1);
          const len = (2.2 + R() * 2).toFixed(1);
          out += `<rect x="${x * CELL}" y="${y * CELL}" width="13" height="13" rx="1.5" fill="${fill}" style="--o:${o};--in:${((1 - d) * 700) | 0}ms;--tw:${tw}s;--len:${len}s"/>`;
        }
      }
    }
    svg.innerHTML = out;
  }
  for (const svg of document.querySelectorAll('.px')) pattern(svg);
  /* light a few squares on every keystroke; a full wave when connected */
  /* effects play on top of each square's breathing (Web Animations), so the build-in never replays */
  const EASE = 'cubic-bezier(.16,1,.3,1)';
  const CEL = { fill: '#5FD6FF', opacity: 1, transform: 'scale(1.35)' };
  function fx(c, frames, ms, delay = 0) {
    if (c.animate && !REDUCED()) c.animate(frames, { duration: ms, delay, easing: EASE });
  }
  const visible = () =>
    [...document.querySelectorAll('.px')].filter((s) => s.getClientRects().length);
  const delayOf = (c) => Number.parseFloat(c.style.getPropertyValue('--in')) || 0;
  /* every square of every visible pattern */
  function eachRect(fn) {
    for (const svg of visible()) for (const c of svg.querySelectorAll('rect')) fn(c);
  }
  /* phone keystrokes: the whole pattern blinks in its own colours */
  function blink() {
    eachRect((c) => fx(c, [{}, { opacity: 1, transform: 'scale(1.14)', offset: 0.25 }, {}], 340));
  }
  /* plus a few celeste squares on the phone field */
  function spark(n) {
    for (const svg of visible()) {
      const r = svg.querySelectorAll('rect');
      for (let i = 0; i < n && r.length; i++) {
        fx(
          r[(Math.random() * r.length) | 0],
          [CEL, { fill: '#5FD6FF', opacity: 1, offset: 0.35 }, {}],
          800
        );
      }
    }
  }
  /* idle pulse every 4s: brightness wave out of the corner, each square keeps its colour */
  setInterval(() => {
    eachRect((c) =>
      fx(c, [{}, { opacity: 1, transform: 'scale(1.2)', offset: 0.3 }, {}], 1000, delayOf(c) * 1.4)
    );
  }, 4200);
  /* closing wave when connected */
  function wave() {
    eachRect((c) => fx(c, [CEL, {}], 1000, delayOf(c)));
  }
  function rule(el, n) {
    if (!el.children.length) {
      for (let i = 0; i < RULE; i++) el.append(document.createElement('i'));
    }
    for (let k = 0; k < RULE; k++) {
      const t = k / (RULE - 1);
      const c = el.children[k];
      c.className = k < n ? 'on' : '';
      c.style.setProperty(
        '--c',
        `rgb(${Math.round(48 + 32 * t)},${Math.round(154 - 100 * t)},${Math.round(231 + 3 * t)})`
      );
    }
  }
  function err(id, on) {
    $(`e-${id}`).classList.toggle('show', on);
    $(id)?.closest('.line')?.classList.toggle('is-error', on);
  }
  function show(id) {
    for (const s of ['s-datos', 's-listo']) {
      $(s).hidden = s !== id;
      $(s).classList.remove('launch', 'land');
    }
  }
  /* the rocket: Byte colours, flame drawn in byte squares */
  const ROCKET =
    '<svg viewBox="0 0 48 92" aria-hidden="true"><path d="M24 2C34 12 38 26 38 44V58H10V44C10 26 14 12 24 2Z" fill="#EDEAF5"/>' +
    '<rect x="10" y="47" width="28" height="3" fill="#309AE7"/><circle cx="24" cy="30" r="5.5" fill="#0A0F1C" stroke="#309AE7" stroke-width="3"/>' +
    '<path d="M10 42L2 58V66L10 60Z M38 42L46 58V66L38 60Z" fill="#5036EA"/><rect x="16" y="58" width="16" height="5" rx="1.5" fill="#3A4358"/>' +
    '<g class="fl"><rect x="17" y="65" width="6" height="6" rx="1" fill="#5FD6FF"/><rect x="25" y="65" width="6" height="6" rx="1" fill="#309AE7"/>' +
    '<rect x="21" y="73" width="6" height="6" rx="1" fill="#5FD6FF"/><rect x="14" y="74" width="5" height="5" rx="1" fill="#5036EA"/>' +
    '<rect x="29" y="74" width="5" height="5" rx="1" fill="#1A6FB8"/><rect x="21" y="82" width="6" height="6" rx="1" fill="#5036EA"/></g></svg>';
  const SMOKE = ['#309AE7', '#5FD6FF', '#5036EA', '#1A6FB8'];
  function puff(x, y) {
    const p = document.createElement('i');
    const z = 5 + Math.random() * 6;
    p.className = 'puff';
    p.style.cssText = `left:${x - z / 2 + (Math.random() - 0.5) * 16}px;top:${y}px;width:${z}px;height:${z}px;background:${SMOKE[(Math.random() * 4) | 0]}`;
    document.body.append(p);
    p.animate(
      [
        { opacity: 1, transform: 'none' },
        { opacity: 0, transform: `translateY(${30 + Math.random() * 70}px) scale(.3)` },
      ],
      { duration: 600 + Math.random() * 350, easing: 'cubic-bezier(.2,.6,.4,1)' }
    ).onfinish = () => p.remove();
  }
  const smooth = (t) => (t < 0.5 ? 2 * t * t : 1 - (2 - 2 * t) ** 2 / 2);
  /* the rocket sits hidden behind the send button, rises out of its top edge nose first,
   then lifts off with its base glued to the button and pulls the whole step away */
  function launch(el, then) {
    if (REDUCED()) return then();
    const r0 = $('b-send').getBoundingClientRect();
    const H = 92;
    const W = 48;
    const k = document.createElement('div');
    k.className = 'rocket';
    k.innerHTML = ROCKET;
    document.body.append(k);
    const legal = el.querySelector('.legal');
    if (legal) {
      legal.style.transition = 'opacity .2s';
      legal.style.opacity = 0;
    } /* clear the rocket's way */
    const kids = [...el.children].filter((c) => !c.classList.contains('corner'));
    const x0 = r0.left + r0.width / 2 - W / 2;
    const VH = innerHeight;
    const t0 = performance.now();
    const IGN = 350;
    const REV = 750;
    const LIFT = 800;
    let last = 0;
    let lifted = false;
    (function frame(now) {
      const t = now - t0;
      let y = 0;
      let jx = 0;
      let blur = 0;
      let out = 0; /* how far the rocket has come out of the button, 0..H */
      if (t < IGN) jx = (Math.random() - 0.5) * 2.4;
      else if (t < IGN + REV) {
        const p = smooth((t - IGN) / REV);
        out = H * p;
        y = -18 * p;
        jx = (Math.random() - 0.5) * 1.6;
      } else {
        const q = Math.min(1, (t - IGN - REV) / LIFT);
        out = H;
        y = -18 - (VH + H) * q * q * q;
        blur = q * 4;
        if (!lifted) {
          lifted = true;
          el.classList.add('launch');
        }
      }
      for (const c of kids) {
        c.style.transform = `translate(${jx}px,${y}px)`;
        c.style.filter = blur ? `blur(${blur.toFixed(1)}px)` : '';
      }
      const line = r0.top + y; /* line: the button's top edge; only what is above it shows */
      const top = line - out;
      k.style.transform = `translate(${x0 + jx}px,${top}px)`;
      k.style.clipPath = `inset(0 0 ${Math.max(0, top + H - line)}px 0)`;
      if (now - last > 28) {
        last = now;
        puff(x0 + W / 2 + jx, Math.min(top + H, line) - 8);
      }
      if (t < IGN + REV + LIFT) requestAnimationFrame(frame);
      else {
        k.remove();
        for (const c of kids) c.style.transform = c.style.filter = '';
        if (legal) legal.style.opacity = '';
        then();
      }
    })(t0);
  }

  /* phone: accept 0412…, 412…, +58 412…, 58412…; keep national digits only */
  function digits(v) {
    let d = v.replace(/\D/g, '');
    const cc = $('cc').value;
    if (d.startsWith(cc) && d.length > (cc === '58' ? 10 : 9)) d = d.slice(cc.length);
    if (cc === '58') return d.replace(/^0+/, '').slice(0, 10);
    return d.slice(0, 14);
  }
  const fmt = (d) =>
    $('cc').value === '58'
      ? [d.slice(0, 3), d.slice(3, 6), d.slice(6)].filter(Boolean).join(' ')
      : d.replace(/(\d{3})(?=\d)/g, '$1 ');
  /* +58: any 10-digit national number; other countries: 7+ digits */
  const telOk = (d) => ($('cc').value === '58' ? d.length === 10 : d.length >= 7);
  function fit() {
    const t = $('tel');
    t.style.fontSize = '';
    for (
      let px = Number.parseFloat(getComputedStyle(t).fontSize);
      t.scrollWidth > t.clientWidth && px > 22;
      px -= 2
    ) {
      t.style.fontSize = `${px}px`;
    }
  }
  function onTel() {
    const d = digits($('tel').value);
    blink();
    spark(4);
    $('tel').value = fmt(d);
    fit();
    rule($('r-tel'), Math.round((d.length / ($('cc').value === '58' ? 10 : 12)) * RULE));
    $('r-tel').classList.remove('is-error');
    err('tel', false);
  }
  /* the country select hugs its current code (mono font, so ch is exact) */
  function ccWidth() {
    $('cc').style.width = `calc(${$('cc').value.length + 1}ch + 16px)`;
  }
  $('tel').addEventListener('input', onTel);
  $('cc').addEventListener('change', () => {
    ccWidth();
    onTel();
  });
  ccWidth();
  $('tel').addEventListener('paste', () => setTimeout(onTel));

  /* cédula: digits with thousands dots */
  $('cedula').addEventListener('input', (e) => {
    const d = e.target.value.replace(/\D/g, '').slice(0, 8); /* max 12.345.678 */
    e.target.value = d.replace(/\B(?=(\d{3})+(?!\d))/g, '.');
    err('cedula', false);
  });
  $('nombre').addEventListener('input', () => err('nombre', false));
  $('empresa').addEventListener('input', () => err('empresa', false));
  /* GO / NEXT: va al siguiente campo sin llenar; si no queda ninguno, oculta el teclado sin enviar.
   El paso siguiente lo da solo el botón */
  const FIELDS = ['nombre', 'cedula', 'empresa', 'tel'];
  for (const [i, id] of FIELDS.entries()) {
    $(id).addEventListener('keydown', (e) => {
      if (e.key !== 'Enter') return;
      e.preventDefault();
      const nxt = FIELDS.slice(i + 1).find((f) => !$(f).value.trim());
      if (nxt) $(nxt).focus();
      else e.target.blur();
    });
  }

  async function post(url, body) {
    if (DEMO) {
      // oxlint-disable-next-line promise/avoid-new -- setTimeout no devuelve promesa: es la latencia falsa del modo demo
      await new Promise((resolve) => setTimeout(resolve, 700));
      return { verified: true };
    }
    const r = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });
    const j = await r.json().catch(() => ({}));
    if (!r.ok) throw j;
    return j;
  }
  function busy(b, on) {
    b.setAttribute('aria-busy', on);
    b.disabled = on;
  }

  $('s-datos').addEventListener('submit', async (e) => {
    e.preventDefault();
    const name = $('nombre').value.trim().replace(/\s+/g, ' ');
    const ced = $('cedula').value.replace(/\D/g, '');
    const company = $('empresa').value.trim().replace(/\s+/g, ' ');
    const tel = digits($('tel').value);
    let bad = null;
    if (name.split(' ').length < 2 || name.length < 4) {
      err('nombre', true);
      bad = bad || 'nombre';
    }
    if (ced.length < 6) {
      err('cedula', true);
      bad = bad || 'cedula';
    }
    if (company.length < 2) {
      err('empresa', true);
      bad = bad || 'empresa';
    }
    if (!telOk(tel)) {
      err('tel', true);
      $('r-tel').classList.add('is-error');
      bad = bad || 'tel';
    }
    if (bad) return $(bad).focus();
    document.activeElement.blur(); /* GO on the phone keyboard: hide it while the rocket flies */
    data = {
      name,
      id_number: document.querySelector('[name=nac]:checked').value + ced,
      company,
      phone: `+${$('cc').value}${tel}`,
    };
    const b = $('b-send');
    busy(b, true);
    err('send', false);
    let r;
    try {
      r = await post(C.send, data);
    } catch {
      busy(b, false);
      err('send', true);
      return;
    }
    busy(b, false);
    launch($('s-datos'), () => done(r?.redirect));
  });

  function done(url) {
    if (url) $('b-go').href = url;
    show('s-listo');
    rule($('r-done'), RULE);
    wave();
    $('net').classList.add('ok');
    $('net').lastChild.textContent = 'Conectado';
  }

  rule($('r-tel'), 0);
  if (DEMO) {
    const d = document.createElement('p');
    d.className = 'demo';
    d.textContent = 'modo demo';
    document.body.append(d);
  }

  /* preview hook for the design-system page: ?step=listo|datos-error */
  const q = location.search.match(/step=([\w-]+)/);
  if (q) {
    $('nombre').value = 'María González';
    $('cedula').value = '18.456.789';
    $('empresa').value = 'Inversiones Demo';
    $('tel').value = '04125551234';
    onTel();
    data = {
      name: 'María González',
      id_number: 'V18456789',
      company: 'Inversiones Demo',
      phone: '+584125551234',
    };
    if (q[1] === 'listo') done();
    if (q[1] === 'datos-error') {
      $('tel').value = '2125551';
      onTel();
      $('s-datos').requestSubmit();
    }
  }
})();
