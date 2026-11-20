/* Backend hooks. Set before this script, e.g. window.BYTE_PORTAL={send:"/api/send",verify:"/api/verify"}.
   The API identifies the device by the request's IP, so the page sends no device id.
   send   POST {name, id_number:"V18456789", phone:"+584125551234"}  -> 2xx {verified?, redirect?}
          verified:true = number already verified on an earlier day: no code, straight to "conectado"
          429 {error:"too_soon"} = a code went out less than 60 s ago and still works: treat as sent
   verify POST {phone, code}                                        -> 2xx {redirect?}
   Without endpoints the portal runs in demo mode (code 123456). */
(function(){
var C=window.BYTE_PORTAL||{},DEMO=!C.send,$=function(i){return document.getElementById(i)};
var RULE=40,WAIT=60,data={},tick;

/* the byte pattern: drawn once per corner, densest at the bottom-right, azul -> violeta */
var PAL=['#161D2E','#3A4358','#1A6FB8','#309AE7','#5FD6FF','#309AE7','#5036EA','#5036EA'];
function pattern(svg){var pal=PAL,seed=11,R=function(){seed=(seed*16807)%2147483647;return seed/2147483647},
    C=15,cols=svg.viewBox.baseVal.width/C|0,rows=svg.viewBox.baseVal.height/C|0,out='';
  for(var y=0;y<rows;y++)for(var x=0;x<cols;x++){var d=x/cols*.55+y/rows*.45;if(R()<d*1.2-.28){
    out+='<rect x="'+x*C+'" y="'+y*C+'" width="13" height="13" rx="1.5" fill="'+pal[Math.min(7,d*8*(.7+R()*.45)|0)]+
    '" style="--o:'+Math.min(1,.3+d*.75).toFixed(2)+';--in:'+((1-d)*700|0)+'ms;--tw:'+(R()*5).toFixed(1)+'s;--len:'+(2.2+R()*2).toFixed(1)+'s"/>'}}
  svg.innerHTML=out}
document.querySelectorAll('.px').forEach(pattern);
/* light a few squares on every keystroke; a full wave when connected */
/* effects play on top of each square's breathing (Web Animations), so the build-in never replays */
var EASE='cubic-bezier(.16,1,.3,1)',CEL={fill:'#5FD6FF',opacity:1,transform:'scale(1.35)'};
function fx(c,frames,ms,delay){if(c.animate&&!matchMedia('(prefers-reduced-motion: reduce)').matches)c.animate(frames,{duration:ms,delay:delay||0,easing:EASE})}
function visible(){return [].filter.call(document.querySelectorAll('.px'),function(s){return s.getClientRects().length})}
function delayOf(c){return parseFloat(c.style.getPropertyValue('--in'))||0}
/* phone and code keystrokes: the whole pattern blinks in its own colours */
function blink(){visible().forEach(function(svg){[].forEach.call(svg.querySelectorAll('rect'),function(c){
  fx(c,[{},{opacity:1,transform:'scale(1.14)',offset:.25},{}],340)})})}
/* plus a few celeste squares on the phone and code fields */
function spark(n){visible().forEach(function(svg){var r=svg.querySelectorAll('rect');
  for(var i=0;i<n&&r.length;i++)fx(r[Math.random()*r.length|0],[CEL,{fill:'#5FD6FF',opacity:1,offset:.35},{}],800)})}
/* idle pulse every 4s: brightness wave out of the corner, each square keeps its colour */
setInterval(function(){visible().forEach(function(svg){[].forEach.call(svg.querySelectorAll('rect'),function(c){
  fx(c,[{},{opacity:1,transform:'scale(1.2)',offset:.3},{}],1000,delayOf(c)*1.4)})})},4200);
/* closing wave when connected */
function wave(){visible().forEach(function(svg){[].forEach.call(svg.querySelectorAll('rect'),function(c){fx(c,[CEL,{}],1000,delayOf(c))})})}
function rule(el,n){if(!el.children.length)for(var i=0;i<RULE;i++)el.appendChild(document.createElement('i'));
  for(var k=0;k<RULE;k++){var t=k/(RULE-1),on=k<n,c=el.children[k];c.className=on?'on':'';
    c.style.setProperty('--c','rgb('+Math.round(48+32*t)+','+Math.round(154-100*t)+','+Math.round(231+3*t)+')')}}
function err(id,on){var e=$('e-'+id);e.classList.toggle('show',on);var f=$(id);if(f&&f.closest('.line'))f.closest('.line').classList.toggle('is-error',on)}
function show(id){['s-datos','s-codigo','s-listo'].forEach(function(s){$(s).hidden=s!==id;$(s).classList.remove('launch','land')})}
/* the rocket: Byte colours, flame drawn in byte squares */
var ROCKET='<svg viewBox="0 0 48 92" aria-hidden="true"><path d="M24 2C34 12 38 26 38 44V58H10V44C10 26 14 12 24 2Z" fill="#EDEAF5"/>'+
  '<rect x="10" y="47" width="28" height="3" fill="#309AE7"/><circle cx="24" cy="30" r="5.5" fill="#0A0F1C" stroke="#309AE7" stroke-width="3"/>'+
  '<path d="M10 42L2 58V66L10 60Z M38 42L46 58V66L38 60Z" fill="#5036EA"/><rect x="16" y="58" width="16" height="5" rx="1.5" fill="#3A4358"/>'+
  '<g class="fl"><rect x="17" y="65" width="6" height="6" rx="1" fill="#5FD6FF"/><rect x="25" y="65" width="6" height="6" rx="1" fill="#309AE7"/>'+
  '<rect x="21" y="73" width="6" height="6" rx="1" fill="#5FD6FF"/><rect x="14" y="74" width="5" height="5" rx="1" fill="#5036EA"/>'+
  '<rect x="29" y="74" width="5" height="5" rx="1" fill="#1A6FB8"/><rect x="21" y="82" width="6" height="6" rx="1" fill="#5036EA"/></g></svg>';
var SMOKE=['#309AE7','#5FD6FF','#5036EA','#1A6FB8'];
function puff(x,y){var p=document.createElement('i'),z=5+Math.random()*6;p.className='puff';
  p.style.cssText='left:'+(x-z/2+(Math.random()-.5)*16)+'px;top:'+y+'px;width:'+z+'px;height:'+z+'px;background:'+SMOKE[Math.random()*4|0];
  document.body.appendChild(p);p.animate([{opacity:1,transform:'none'},{opacity:0,transform:'translateY('+(30+Math.random()*70)+'px) scale(.3)'}],
  {duration:600+Math.random()*350,easing:'cubic-bezier(.2,.6,.4,1)'}).onfinish=function(){p.remove()}}
/* the rocket sits hidden behind the send button, rises out of its top edge nose first,
   then lifts off with its base glued to the button and pulls the whole step away */
function launch(el,then){if(matchMedia('(prefers-reduced-motion: reduce)').matches)return then();
  var b=$('b-send'),r0=b.getBoundingClientRect(),H=92,W=48,k=document.createElement('div');
  k.className='rocket';k.innerHTML=ROCKET;document.body.appendChild(k);
  var legal=el.querySelector('.legal');if(legal){legal.style.transition='opacity .2s';legal.style.opacity=0}  /* clear the rocket's way */
  var kids=[].filter.call(el.children,function(c){return !c.classList.contains('corner')}),
    x0=r0.left+r0.width/2-W/2,VH=innerHeight,t0=performance.now(),last=0,lifted=false,
    IGN=350,REV=750,LIFT=800;
  function smooth(t){return t<.5?2*t*t:1-Math.pow(2-2*t,2)/2}
  (function frame(now){var t=now-t0,y=0,jx=0,blur=0;
    var out=0;  /* how far the rocket has come out of the button, 0..H */
    if(t<IGN)jx=(Math.random()-.5)*2.4;
    else if(t<IGN+REV){var p=smooth((t-IGN)/REV);out=H*p;y=-18*p;jx=(Math.random()-.5)*1.6}
    else{var q=Math.min(1,(t-IGN-REV)/LIFT);out=H;y=-18-(VH+H)*q*q*q;blur=q*4;
      if(!lifted){lifted=true;el.classList.add('launch')}}
    kids.forEach(function(c){c.style.transform='translate('+jx+'px,'+y+'px)';c.style.filter=blur?'blur('+blur.toFixed(1)+'px)':''});
    var line=r0.top+y,top=line-out;  /* line: the button's top edge; only what is above it shows */
    k.style.transform='translate('+(x0+jx)+'px,'+top+'px)';
    k.style.clipPath='inset(0 0 '+Math.max(0,top+H-line)+'px 0)';
    if(now-last>28){last=now;puff(x0+W/2+jx,Math.min(top+H,line)-8)}
    if(t<IGN+REV+LIFT)requestAnimationFrame(frame);
    else{k.remove();kids.forEach(function(c){c.style.transform=c.style.filter=''});if(legal)legal.style.opacity='';then()}})(t0)}

/* phone: accept 0412…, 412…, +58 412…, 58412…; keep national digits only */
function digits(v){var d=v.replace(/\D/g,''),cc=$('cc').value;
  if(d.indexOf(cc)===0&&d.length>(cc==='58'?10:9))d=d.slice(cc.length);
  if(cc==='58'){d=d.replace(/^0+/,'');return d.slice(0,10)}
  return d.slice(0,14)}
function fmt(d){return $('cc').value==='58'?[d.slice(0,3),d.slice(3,6),d.slice(6)].filter(Boolean).join(' '):d.replace(/(\d{3})(?=\d)/g,'$1 ')}
/* +58: any 10-digit national number; other countries: 7+ digits */
function telOk(d){return $('cc').value==='58'?d.length===10:d.length>=7}
function fit(){var t=$('tel');t.style.fontSize='';for(var px=parseFloat(getComputedStyle(t).fontSize);t.scrollWidth>t.clientWidth&&px>22;px-=2)t.style.fontSize=px+'px'}
function onTel(){var d=digits($('tel').value);blink();spark(4);$('tel').value=fmt(d);fit();rule($('r-tel'),Math.round(d.length/($('cc').value==='58'?10:12)*RULE));
  $('r-tel').classList.remove('is-error');err('tel',false)}
/* the country select hugs its current code (mono font, so ch is exact) */
function ccWidth(){$('cc').style.width='calc('+($('cc').value.length+1)+'ch + 16px)'}
$('tel').addEventListener('input',onTel);$('cc').addEventListener('change',function(){ccWidth();onTel()});ccWidth();
$('tel').addEventListener('paste',function(){setTimeout(onTel)});

/* cédula: digits with thousands dots */
$('cedula').addEventListener('input',function(){var d=this.value.replace(/\D/g,'').slice(0,8);  /* max 12.345.678 */
  this.value=d.replace(/\B(?=(\d{3})+(?!\d))/g,'.');err('cedula',false)});
$('nombre').addEventListener('input',function(){err('nombre',false)});
/* GO / NEXT: va al siguiente campo sin llenar; si no queda ninguno, oculta el teclado sin enviar.
   El paso siguiente lo da solo el botón */
['nombre','cedula','tel'].forEach(function(id,i,all){$(id).addEventListener('keydown',function(e){if(e.key!=='Enter')return;e.preventDefault();
  var nxt=all.slice(i+1).filter(function(f){return !$(f).value.trim()})[0];
  if(nxt)$(nxt).focus();else this.blur()})});
$('code').addEventListener('keydown',function(e){if(e.key==='Enter'){e.preventDefault();this.blur()}});

function post(url,body){return DEMO?new Promise(function(ok){setTimeout(ok,700)}):fetch(url,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)}).then(function(r){return r.json().catch(function(){return{}}).then(function(j){if(!r.ok)throw j;return j})})}
function busy(b,on){b.setAttribute('aria-busy',on);b.disabled=on}

$('s-datos').addEventListener('submit',function(e){e.preventDefault();
  var n=$('nombre').value.trim().replace(/\s+/g,' '),c=$('cedula').value.replace(/\D/g,''),t=digits($('tel').value),bad=null;
  if(n.split(' ').length<2||n.length<4){err('nombre',true);bad=bad||'nombre'}
  if(c.length<6){err('cedula',true);bad=bad||'cedula'}
  if(!telOk(t)){err('tel',true);$('r-tel').classList.add('is-error');bad=bad||'tel'}
  if(bad)return $(bad).focus();
  document.activeElement.blur();  /* GO on the phone keyboard: hide it while the rocket flies */
  data={name:n,id_number:document.querySelector('[name=nac]:checked').value+c,phone:'+'+$('cc').value+t};
  var b=$('b-send');busy(b,true);err('send',false);
  function next(r){busy(b,false);launch($('s-datos'),r&&r.verified?function(){done(r.redirect)}:function(){toCode();$('s-codigo').classList.add('land')})}
  /* too_soon: el código salió hace menos de 60 s y sigue valiendo, así que pasa al paso del código */
  post(C.send,data).then(next,function(e){if(e&&e.error==='too_soon')return next({});busy(b,false);err('send',true)})});

function toCode(){$('sent-to').textContent='+'+$('cc').value+' '+fmt(digits($('tel').value));show('s-codigo');paint();timer();$('code').focus()}
/* the pixel rule drains with the resend countdown */
function timer(){var left=WAIT;clearInterval(tick);$('b-resend').disabled=true;
  function step(){rule($('r-time'),Math.round(left/WAIT*RULE));
    $('t-left').textContent=left>0?'Reenviar en '+Math.floor(left/60)+':'+String(left%60).padStart(2,'0'):'¿No te llegó?';
    $('b-resend').hidden=left>0;if(left<=0){clearInterval(tick);$('b-resend').disabled=false}left--}
  step();tick=setInterval(step,1000)}
$('b-resend').addEventListener('click',function(){var b=this;b.disabled=true;post(C.send,data).then(timer,function(e){if(e&&e.error==='too_soon')return timer();b.disabled=false})});
$('b-edit').addEventListener('click',function(){clearInterval(tick);show('s-datos');$('tel').focus()});

function paint(){var v=$('code').value,cells=$('otp').querySelectorAll('span'),f=document.activeElement===$('code');
  cells.forEach(function(s,i){s.textContent=v[i]||'';s.classList.toggle('cur',f&&i===Math.min(v.length,5))})}
$('code').addEventListener('input',function(){this.value=this.value.replace(/\D/g,'').slice(0,6);blink();spark(6);$('otp').classList.remove('is-error');$('e-code').classList.remove('show');paint()});
$('code').addEventListener('focus',paint);$('code').addEventListener('blur',paint);

$('s-codigo').addEventListener('submit',function(e){e.preventDefault();var v=$('code').value,b=$('b-verify');
  if(v.length<6)return $('code').focus();busy(b,true);
  var p=DEMO?new Promise(function(ok,no){setTimeout(v==='123456'?ok:no,700)}):post(C.verify,{phone:data.phone,code:v});
  p.then(function(r){busy(b,false);clearInterval(tick);done(r&&r.redirect)},function(){busy(b,false);$('otp').classList.add('is-error');$('e-code').classList.add('show');$('code').select()})});

function done(url){if(url)$('b-go').href=url;show('s-listo');rule($('r-done'),RULE);wave();
  $('net').classList.add('ok');$('net').lastChild.textContent='Conectado'}

rule($('r-tel'),0);
if(DEMO){var d=document.createElement('p');d.className='demo';d.textContent='modo demo · código 123456';document.body.appendChild(d)}

/* preview hook for the design-system page: ?step=codigo|error|listo */
var q=location.search.match(/step=([\w-]+)/);
if(q){$('nombre').value='María González';$('cedula').value='18.456.789';$('tel').value='04125551234';onTel();
  data={name:'María González',id_number:'V18456789',phone:'+584125551234'};
  if(q[1]==='codigo'||q[1]==='error'){toCode();if(q[1]==='error'){$('code').value='483920';paint();$('otp').classList.add('is-error');$('e-code').classList.add('show')}}
  if(q[1]==='listo')done();
  if(q[1]==='datos-error'){$('tel').value='2125551';onTel();$('s-datos').requestSubmit()}}
})();
