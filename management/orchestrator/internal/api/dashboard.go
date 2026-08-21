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
:root{--bg:#0f1115;--card:#1a1d24;--text:#e8eaed;--muted:#9aa0a6;--accent:#7c9cff;--good:#3dd68c;--warn:#f5a524;--bad:#f76e6e;--border:#2a2f3a}
*{box-sizing:border-box}body{margin:0;font-family:system-ui,sans-serif;background:var(--bg);color:var(--text)}
header{padding:1rem 1.5rem;border-bottom:1px solid var(--border);display:flex;justify-content:space-between;flex-wrap:wrap;gap:.75rem;align-items:center}
h1{margin:0;font-size:1.2rem}.sub{color:var(--muted);font-size:.85rem}
nav{display:flex;gap:.35rem;flex-wrap:wrap}
nav button{background:transparent;color:var(--muted);border:1px solid transparent;border-radius:8px;padding:.4rem .75rem;cursor:pointer;font-size:.85rem}
nav button.active,nav button:hover{color:var(--text);background:var(--card);border-color:var(--border)}
main{padding:1.25rem 1.5rem;max-width:1200px;margin:0 auto}
.panel{display:none}.panel.active{display:block}
.strip{display:grid;grid-template-columns:repeat(auto-fit,minmax(120px,1fr));gap:.75rem;margin-bottom:1.25rem}
.stat{background:var(--card);border-radius:10px;padding:.9rem}
.stat .label{color:var(--muted);font-size:.7rem;text-transform:uppercase}
.stat .value{font-size:1.35rem;font-weight:600;margin-top:.2rem}
.actions{margin-bottom:1rem;display:flex;gap:.5rem;flex-wrap:wrap;align-items:center}
.btn{background:var(--accent);color:#0f1115;border:none;border-radius:8px;padding:.45rem .85rem;font-weight:600;cursor:pointer;font-size:.85rem}
.btn.secondary{background:#2a2f3a;color:var(--text)}.btn.danger{background:#5c2a2a;color:var(--bad)}
table{width:100%;border-collapse:collapse;background:var(--card);border-radius:10px;overflow:hidden;margin-bottom:1.5rem}
th,td{text-align:left;padding:.65rem .85rem;border-bottom:1px solid var(--border);font-size:.85rem}
th{color:var(--muted);font-size:.7rem;text-transform:uppercase}
.pill{display:inline-block;padding:.12rem .45rem;border-radius:999px;font-size:.72rem}
.thriving,.running{background:#1a3d2e;color:var(--good)}.stable{background:#1e2a3d;color:var(--accent)}
.under_pressure{background:#3d2e1a;color:var(--warn)}.critical,.collapsed{background:#3d1a1a;color:var(--bad)}
.stopped{background:#2a2f3a;color:var(--muted)}.empty{color:var(--muted);padding:1.5rem;text-align:center}
code{font-size:.78rem;color:var(--muted)}.row-actions{display:flex;gap:.35rem;flex-wrap:wrap}
.msg{color:var(--muted);font-size:.8rem}h2{font-size:1rem;margin:0 0 .75rem}
.note{background:#1a1d24;border:1px solid var(--border);border-radius:10px;padding:1rem;color:var(--muted);font-size:.9rem;margin-bottom:1rem;line-height:1.45}
.screens{display:grid;grid-template-columns:repeat(auto-fill,minmax(180px,1fr));gap:1rem}
.phone{background:var(--card);border-radius:12px;overflow:hidden;border:1px solid var(--border)}
.phone .meta{padding:.5rem .65rem;font-size:.75rem;color:var(--muted);display:flex;justify-content:space-between;gap:.5rem}
.phone img{width:100%;display:block;background:#000;min-height:120px;object-fit:contain;aspect-ratio:9/16}
.phone .ph{padding:2rem .5rem;text-align:center;color:var(--muted);font-size:.8rem;aspect-ratio:9/16;display:flex;align-items:center;justify-content:center}
</style>
</head>
<body>
<header>
<div><h1>Unborn</h1><div class="sub">Farming with souls · <span id="lic-line">license…</span></div></div>
<nav>
<button type="button" class="active" data-p="screens">Screens</button>
<button type="button" data-p="pop">Population</button>
<button type="button" data-p="bodies">Bodies</button>
<button type="button" data-p="pb">Playbooks</button>
<button type="button" data-p="px">Proxies</button>
</nav>
</header>
<main>
<div class="strip">
<div class="stat"><div class="label">Personas</div><div class="value" id="s-p">—</div></div>
<div class="stat"><div class="label">Bodies</div><div class="value" id="s-b">—</div></div>
<div class="stat"><div class="label">Real screens</div><div class="value" id="s-r">—</div></div>
<div class="stat"><div class="label">License max</div><div class="value" id="s-l">—</div></div>
</div>

<section id="panel-screens" class="panel active">
<div class="note">
<strong>Screens need a real Redroid body.</strong> Simulated bodies have no Android UI — only data/API.
Start a persona → <em>Real body</em> (host needs binder modules + Docker + <code>adb</code>). Then frames appear here via ADB screenshot (refreshes every few seconds).
</div>
<div class="actions"><button type="button" class="btn" id="ref-sc">Refresh frames</button><span class="msg" id="msg-sc"></span></div>
<div class="screens" id="grid"></div>
</section>

<section id="panel-pop" class="panel">
<div class="actions"><button type="button" class="btn" id="ref">Refresh</button>
<button type="button" class="btn secondary" id="mk">+ Persona</button><span class="msg" id="msg"></span></div>
<table><thead><tr><th>Name</th><th>Location</th><th>Engagement</th><th>Vitality</th><th>Body</th><th></th></tr></thead>
<tbody id="tp"></tbody></table>
</section>

<section id="panel-bodies" class="panel">
<h2>Bodies</h2>
<table><thead><tr><th>ID</th><th>Persona</th><th>State</th><th>Sim</th><th>ADB</th><th></th></tr></thead>
<tbody id="tb"></tbody></table>
</section>

<section id="panel-pb" class="panel">
<h2>Playbooks</h2>
<table><thead><tr><th>Name</th><th>Kind</th><th>Description</th></tr></thead><tbody id="tpb"></tbody></table>
</section>

<section id="panel-px" class="panel">
<h2>Proxies</h2>
<table><thead><tr><th>Persona</th><th>Host</th><th>Port</th><th>Type</th></tr></thead><tbody id="tpx"></tbody></table>
</section>
</main>
<script>
const $=id=>document.getElementById(id);
let S={personas:[],instances:[],vitality:[],playbooks:[],proxies:[],license:null};
document.querySelectorAll('nav button').forEach(b=>{
b.onclick=()=>{
document.querySelectorAll('nav button').forEach(x=>x.classList.remove('active'));
document.querySelectorAll('.panel').forEach(x=>x.classList.remove('active'));
b.classList.add('active');$('panel-'+b.dataset.p).classList.add('active');
if(b.dataset.p==='screens') refreshScreens();
};
});
function lv(v){return v>=80?'thriving':v>=55?'stable':v>=30?'under_pressure':v>=10?'critical':'collapsed'}
async function load(){
const [personas,vitality,instances,license,playbooks,proxies]=await Promise.all([
fetch('/v1/personas').then(r=>r.json()),fetch('/v1/vitality').then(r=>r.json()),
fetch('/v1/instances').then(r=>r.json()),fetch('/v1/license').then(r=>r.json()),
fetch('/v1/playbooks').then(r=>r.json()),fetch('/v1/proxies').then(r=>r.json())]);
S={personas:personas||[],instances:instances||[],vitality:vitality||[],playbooks:playbooks||[],proxies:proxies||[],license};
render();refreshScreens();
}
function render(){
const vm={};S.vitality.forEach(v=>vm[v.persona_id]=v);
const bm={};S.instances.forEach(i=>bm[i.persona_id]=i);
const real=S.instances.filter(i=>!i.simulated&&i.state==='running').length;
$('s-p').textContent=S.personas.length;$('s-b').textContent=S.instances.length;$('s-r').textContent=real;
$('s-l').textContent=S.license?S.license.max_instances:'—';
if(S.license)$('lic-line').textContent=(S.license.valid?'valid':'invalid')+' · '+(S.license.tier||'')+' · max '+S.license.max_instances;
$('tp').innerHTML=S.personas.length?S.personas.map(p=>{
const v=vm[p.id],val=v?v.value.toFixed(1):'75.0',l=v?lv(v.value):'stable',b=bm[p.id];
const bl=b?(b.simulated?'sim':'real')+' '+(b.state||''):'—';
return '<tr><td>'+(p.display_name||'—')+'</td><td>'+(p.demographics&&p.demographics.location||'—')+'</td><td>'+(p.engagement&&p.engagement.type||'—')+
'</td><td><span class="pill '+l+'">'+val+'</span></td><td><code>'+bl+'</code></td><td class="row-actions">'+
'<button type="button" class="btn secondary" data-s="'+p.id+'">Sim</button>'+
'<button type="button" class="btn secondary" data-r="'+p.id+'">Real</button></td></tr>';
}).join(''):'<tr><td colspan="6" class="empty">No personas</td></tr>';
$('tp').querySelectorAll('[data-s]').forEach(b=>b.onclick=()=>start(b.dataset.s,true));
$('tp').querySelectorAll('[data-r]').forEach(b=>b.onclick=()=>start(b.dataset.r,false));
$('tb').innerHTML=S.instances.length?S.instances.map(i=>{
const st=i.state==='running'?'running':'stopped';
return '<tr><td><code>'+(i.id||'').slice(0,8)+'…</code></td><td><code>'+(i.persona_id||'').slice(0,8)+'…</code></td>'+
'<td><span class="pill '+st+'">'+(i.state||'')+'</span></td><td>'+(i.simulated?'yes':'no')+'</td><td>'+(i.adb_port||'—')+'</td>'+
'<td class="row-actions">'+
'<button type="button" class="btn secondary" data-h="'+i.id+'">Health</button>'+
'<button type="button" class="btn danger" data-x="'+i.id+'">Stop</button></td></tr>';
}).join(''):'<tr><td colspan="6" class="empty">No bodies</td></tr>';
$('tb').querySelectorAll('[data-x]').forEach(b=>b.onclick=async()=>{await fetch('/v1/instances/'+b.dataset.x+'/stop',{method:'POST'});load()});
$('tb').querySelectorAll('[data-h]').forEach(b=>b.onclick=async()=>{const r=await fetch('/v1/instances/'+b.dataset.h+'/health').then(x=>x.json());alert((r.healthy?'OK':'FAIL')+': '+r.reason)});
$('tpb').innerHTML=S.playbooks.length?S.playbooks.map(p=>'<tr><td>'+p.name+'</td><td>'+p.kind+'</td><td>'+(p.description||'')+'</td></tr>').join(''):'<tr><td colspan="3" class="empty">None</td></tr>';
$('tpx').innerHTML=S.proxies.length?S.proxies.map(p=>'<tr><td><code>'+p.persona_id.slice(0,8)+'…</code></td><td>'+p.host+'</td><td>'+p.port+'</td><td>'+(p.type||'')+'</td></tr>').join(''):'<tr><td colspan="4" class="empty">None</td></tr>';
}
function refreshScreens(){
const grid=$('grid');
const real=S.instances.filter(i=>!i.simulated&&i.state==='running');
const sim=S.instances.filter(i=>i.simulated&&i.state==='running');
if(!real.length&&!sim.length){
grid.innerHTML='<div class="empty">No running bodies. Create a persona → start <strong>Real</strong> body to see a screen.</div>';
return;
}
let html='';
real.forEach(i=>{
const ts=Date.now();
html+='<div class="phone"><div class="meta"><span>real · adb '+(i.adb_port||'?')+'</span><code>'+i.id.slice(0,6)+'</code></div>'+
'<img src="/v1/instances/'+i.id+'/screenshot?t='+ts+'" alt="screen" onerror="this.parentElement.querySelector(\'.fail\').style.display=\'flex\';this.style.display=\'none\'"/>'+
'<div class="ph fail" style="display:none">Waiting for boot / adb…</div></div>';
});
sim.forEach(i=>{
html+='<div class="phone"><div class="meta"><span>simulated</span><code>'+i.id.slice(0,6)+'</code></div>'+
'<div class="ph">No screen — simulated body is API-only. Use Real body for Android UI.</div></div>';
});
grid.innerHTML=html;
$('msg-sc').textContent=real.length?'Loaded '+real.length+' real frame slot(s)':'Only simulated — no pixels yet';
}
async function start(id,sim){
$('msg').textContent='Starting…';
const url=sim?'/v1/instances':'/v1/instances?real=true';
const r=await fetch(url,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({persona_id:id,simulated:sim})});
$('msg').textContent=r.ok?'Started':'Fail: '+await r.text();load();
}
$('ref').onclick=()=>load();
$('ref-sc').onclick=()=>refreshScreens();
$('mk').onclick=async()=>{await fetch('/v1/personas',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({display_name:'Persona '+Math.floor(Math.random()*1000),location:'Berlin',timezone:'Europe/Berlin',age_min:25,age_max:30,engagement:'thoughtful_commenter'})});load()};
load();setInterval(load,20000);setInterval(refreshScreens,5000);
</script>
</body>
</html>`

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(dashboardHTML))
}
