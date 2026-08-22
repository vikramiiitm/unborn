package api

import (
	"net/http"
)

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8"/><meta name="viewport" content="width=device-width,initial-scale=1"/>
<title>Unborn</title>
<style>
:root{--bg:#0f1115;--card:#1a1d24;--text:#e8eaed;--muted:#9aa0a6;--accent:#7c9cff;--good:#3dd68c;--bad:#f76e6e;--border:#2a2f3a}
*{box-sizing:border-box}body{margin:0;font-family:system-ui,sans-serif;background:var(--bg);color:var(--text)}
header{padding:1rem 1.5rem;border-bottom:1px solid var(--border);display:flex;justify-content:space-between;flex-wrap:wrap;gap:.75rem;align-items:center}
h1{margin:0;font-size:1.2rem}.sub{color:var(--muted);font-size:.85rem}
nav{display:flex;gap:.35rem;flex-wrap:wrap}
nav button{background:transparent;color:var(--muted);border:1px solid transparent;border-radius:8px;padding:.4rem .75rem;cursor:pointer;font-size:.85rem}
nav button.active,nav button:hover{color:var(--text);background:var(--card);border-color:var(--border)}
main{padding:1.25rem 1.5rem;max-width:1280px;margin:0 auto}
.panel{display:none}.panel.active{display:block}
.strip{display:grid;grid-template-columns:repeat(auto-fit,minmax(110px,1fr));gap:.75rem;margin-bottom:1.25rem}
.stat{background:var(--card);border-radius:10px;padding:.85rem}
.stat .label{color:var(--muted);font-size:.7rem;text-transform:uppercase}
.stat .value{font-size:1.3rem;font-weight:600;margin-top:.15rem}
.actions{margin-bottom:1rem;display:flex;gap:.5rem;flex-wrap:wrap;align-items:center}
.btn{background:var(--accent);color:#0f1115;border:none;border-radius:8px;padding:.4rem .75rem;font-weight:600;cursor:pointer;font-size:.8rem}
.btn.secondary{background:#2a2f3a;color:var(--text)}.btn.danger{background:#5c2a2a;color:var(--bad)}
table{width:100%;border-collapse:collapse;background:var(--card);border-radius:10px;overflow:hidden;margin-bottom:1.25rem}
th,td{text-align:left;padding:.6rem .8rem;border-bottom:1px solid var(--border);font-size:.85rem}
th{color:var(--muted);font-size:.7rem;text-transform:uppercase}
.pill{display:inline-block;padding:.1rem .4rem;border-radius:999px;font-size:.7rem}
.running{background:#1a3d2e;color:var(--good)}.stopped{background:#2a2f3a;color:var(--muted)}
.empty{color:var(--muted);padding:1.5rem;text-align:center}code{font-size:.75rem;color:var(--muted)}
.row-actions{display:flex;gap:.3rem;flex-wrap:wrap}.msg{color:var(--muted);font-size:.8rem}
.note{background:var(--card);border:1px solid var(--border);border-radius:10px;padding:.85rem;color:var(--muted);font-size:.85rem;margin-bottom:1rem;line-height:1.4}
.screens{display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:1rem}
.phone{background:var(--card);border-radius:14px;overflow:hidden;border:1px solid var(--border);display:flex;flex-direction:column}
.phone .meta{padding:.45rem .6rem;font-size:.72rem;color:var(--muted);display:flex;justify-content:space-between}
.phone .stage{position:relative;background:#000;cursor:crosshair;user-select:none}
.phone img{width:100%;display:block;aspect-ratio:9/16;object-fit:contain}
.phone .ph{padding:2rem .5rem;text-align:center;color:var(--muted);font-size:.8rem;aspect-ratio:9/16;display:flex;align-items:center;justify-content:center}
.phone .ctrls{display:flex;gap:.3rem;padding:.45rem;flex-wrap:wrap;border-top:1px solid var(--border)}
.phone .ctrls button{flex:1;min-width:3rem}
.hint{font-size:.75rem;color:var(--muted);padding:0 .6rem .5rem}
</style>
</head>
<body>
<header>
<div><h1>Unborn</h1><div class="sub">Interactive · <span id="lic-line">…</span></div></div>
<nav>
<button type="button" class="active" data-p="screens">Screens</button>
<button type="button" data-p="pop">Population</button>
<button type="button" data-p="bodies">Bodies</button>
<button type="button" data-p="pb">Playbooks</button>
</nav>
</header>
<main>
<div class="strip">
<div class="stat"><div class="label">Personas</div><div class="value" id="s-p">—</div></div>
<div class="stat"><div class="label">Bodies</div><div class="value" id="s-b">—</div></div>
<div class="stat"><div class="label">Real</div><div class="value" id="s-r">—</div></div>
<div class="stat"><div class="label">License</div><div class="value" id="s-l">—</div></div>
</div>
<section id="panel-screens" class="panel active">
<div class="note"><strong>Click</strong> = tap · <strong>drag</strong> = swipe · Home / Back / Recents. Coords map to 1080×1920.</div>
<div class="actions"><button type="button" class="btn" id="ref-sc">Refresh</button><span class="msg" id="msg-sc"></span></div>
<div class="screens" id="grid"></div>
</section>
<section id="panel-pop" class="panel">
<div class="actions"><button type="button" class="btn" id="ref">Refresh</button>
<button type="button" class="btn secondary" id="mk">+ Persona</button><span class="msg" id="msg"></span></div>
<table><thead><tr><th>Name</th><th>Location</th><th>Body</th><th></th></tr></thead><tbody id="tp"></tbody></table>
</section>
<section id="panel-bodies" class="panel">
<table><thead><tr><th>ID</th><th>State</th><th>Sim</th><th>ADB</th><th></th></tr></thead><tbody id="tb"></tbody></table>
</section>
<section id="panel-pb" class="panel">
<table><thead><tr><th>Name</th><th>Kind</th></tr></thead><tbody id="tpb"></tbody></table>
</section>
</main>
<script>
const $=id=>document.getElementById(id); const DW=1080, DH=1920;
let S={personas:[],instances:[],license:null,playbooks:[]};
document.querySelectorAll('nav button').forEach(b=>{
 b.onclick=()=>{document.querySelectorAll('nav button').forEach(x=>x.classList.remove('active'));
 document.querySelectorAll('.panel').forEach(x=>x.classList.remove('active'));
 b.classList.add('active');$('panel-'+b.dataset.p).classList.add('active');
 if(b.dataset.p==='screens') refreshScreens();};
});
async function load(){
 const [personas,instances,license,playbooks]=await Promise.all([
  fetch('/v1/personas').then(r=>r.json()),fetch('/v1/instances').then(r=>r.json()),
  fetch('/v1/license').then(r=>r.json()),fetch('/v1/playbooks').then(r=>r.json())]);
 S={personas:personas||[],instances:instances||[],license,playbooks:playbooks||[]}; render(); refreshScreens();
}
function render(){
 const bm={}; S.instances.forEach(i=>bm[i.persona_id]=i);
 const real=S.instances.filter(i=>!i.simulated&&i.state==='running').length;
 $('s-p').textContent=S.personas.length; $('s-b').textContent=S.instances.length; $('s-r').textContent=real;
 $('s-l').textContent=S.license?S.license.max_instances:'—';
 if(S.license)$('lic-line').textContent=(S.license.valid?'valid':'invalid')+' · max '+S.license.max_instances;
 $('tp').innerHTML=S.personas.length?S.personas.map(p=>{
  const b=bm[p.id]; const bl=b?(b.simulated?'sim':'real')+' '+(b.state||''):'—';
  return '<tr><td>'+(p.display_name||'—')+'</td><td>'+(p.demographics&&p.demographics.location||'—')+
   '</td><td><code>'+bl+'</code></td><td class="row-actions">'+
   '<button type="button" class="btn secondary" data-s="'+p.id+'">Sim</button>'+
   '<button type="button" class="btn secondary" data-r="'+p.id+'">Real</button></td></tr>';
 }).join(''):'<tr><td colspan="4" class="empty">No personas</td></tr>';
 $('tp').querySelectorAll('[data-s]').forEach(b=>b.onclick=()=>start(b.dataset.s,true));
 $('tp').querySelectorAll('[data-r]').forEach(b=>b.onclick=()=>start(b.dataset.r,false));
 $('tb').innerHTML=S.instances.length?S.instances.map(i=>{
  const st=i.state==='running'?'running':'stopped';
  return '<tr><td><code>'+i.id.slice(0,8)+'</code></td><td><span class="pill '+st+'">'+(i.state||'')+
   '</span></td><td>'+(i.simulated?'yes':'no')+'</td><td>'+(i.adb_port||'—')+
   '</td><td><button type="button" class="btn danger" data-x="'+i.id+'">Stop</button></td></tr>';
 }).join(''):'<tr><td colspan="5" class="empty">No bodies</td></tr>';
 $('tb').querySelectorAll('[data-x]').forEach(b=>b.onclick=async()=>{await fetch('/v1/instances/'+b.dataset.x+'/stop',{method:'POST'});load()});
 $('tpb').innerHTML=S.playbooks.length?S.playbooks.map(p=>'<tr><td>'+p.name+'</td><td>'+p.kind+'</td></tr>').join(''):'<tr><td colspan="2" class="empty">None</td></tr>';
}
function deviceXY(img,cx,cy){const r=img.getBoundingClientRect();
 return {x:Math.max(0,Math.min(DW-1,Math.round((cx-r.left)/r.width*DW))),
 y:Math.max(0,Math.min(DH-1,Math.round((cy-r.top)/r.height*DH)))};}
async function postInput(id,path,body){
 const r=await fetch('/v1/instances/'+id+'/input/'+path,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
 $('msg-sc').textContent=r.ok?path+' ok':path+' fail'; setTimeout(refreshScreens,350);
}
function wirePhone(el,id){
 const stage=el.querySelector('.stage'), img=el.querySelector('img'); if(!stage||!img) return;
 let down=null;
 stage.onpointerdown=e=>{e.preventDefault(); stage.setPointerCapture(e.pointerId); down=deviceXY(img,e.clientX,e.clientY);};
 stage.onpointerup=e=>{if(!down)return; const up=deviceXY(img,e.clientX,e.clientY);
  if(Math.abs(up.x-down.x)+Math.abs(up.y-down.y)<20) postInput(id,'tap',{x:up.x,y:up.y});
  else postInput(id,'swipe',{x1:down.x,y1:down.y,x2:up.x,y2:up.y,ms:280}); down=null;};
 el.querySelectorAll('[data-key]').forEach(b=>b.onclick=()=>postInput(id,'key',{keycode:+b.dataset.key}));
}
function refreshScreens(){
 const grid=$('grid');
 const real=S.instances.filter(i=>!i.simulated&&i.state==='running');
 const sim=S.instances.filter(i=>i.simulated&&i.state==='running');
 if(!real.length&&!sim.length){grid.innerHTML='<div class="empty">Start a Real body to control a screen.</div>';return;}
 let html='';
 real.forEach(i=>{html+='<div class="phone" data-id="'+i.id+'"><div class="meta"><span>real · '+(i.adb_port||'?')+'</span><code>'+i.id.slice(0,6)+'</code></div>'+
  '<div class="stage"><img src="/v1/instances/'+i.id+'/screenshot?t='+Date.now()+'" draggable="false" onerror="this.style.opacity=.3"/></div>'+
  '<div class="ctrls"><button type="button" class="btn secondary" data-key="4">Back</button>'+
  '<button type="button" class="btn secondary" data-key="3">Home</button>'+
  '<button type="button" class="btn secondary" data-key="187">Recents</button></div>'+
  '<div class="hint">Click = tap · drag = swipe</div></div>';});
 sim.forEach(i=>{html+='<div class="phone"><div class="meta">sim</div><div class="ph">API only — use Real</div></div>';});
 grid.innerHTML=html;
 grid.querySelectorAll('.phone[data-id]').forEach(el=>wirePhone(el,el.dataset.id));
 $('msg-sc').textContent=real.length?real.length+' interactive':'sim only';
}
async function start(id,sim){
 $('msg').textContent='Starting…';
 const url=sim?'/v1/instances':'/v1/instances?real=true';
 const r=await fetch(url,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({persona_id:id,simulated:sim})});
 $('msg').textContent=r.ok?'Started':'Fail'; load();
}
$('ref').onclick=()=>load(); $('ref-sc').onclick=()=>refreshScreens();
$('mk').onclick=async()=>{await fetch('/v1/personas',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({display_name:'Persona '+Math.floor(Math.random()*1000),location:'Berlin',timezone:'Europe/Berlin'})});load()};
load(); setInterval(load,25000); setInterval(refreshScreens,4000);
</script>
</body>
</html>`

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(dashboardHTML))
}
