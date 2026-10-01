(async () => {
  const gate=document.querySelector('#static-gate');
  const shell=document.querySelector('#site-shell');
  if(!gate||!shell)return;
  const prefix=location.pathname.includes('/projects/')?'../':'';
  let cfg={gateEnabled:false};
  try{cfg=await (await fetch(prefix+'site-config.json',{cache:'no-store'})).json();}catch(e){console.error(e)}
  if(!cfg.gateEnabled){shell.classList.remove('locked');return}
  const key='mpr-static-unlocked:'+cfg.passwordSHA256;
  if(sessionStorage.getItem(key)==='1'){shell.classList.remove('locked');return}
  gate.classList.remove('hidden');
  const form=document.querySelector('#gate-form'); const input=document.querySelector('#gate-password'); const err=document.querySelector('#gate-error');
  form.addEventListener('submit',async(ev)=>{ev.preventDefault(); const bytes=new TextEncoder().encode(input.value); const digest=await crypto.subtle.digest('SHA-256',bytes); const hex=[...new Uint8Array(digest)].map(b=>b.toString(16).padStart(2,'0')).join(''); if(hex===cfg.passwordSHA256){sessionStorage.setItem(key,'1');gate.classList.add('hidden');shell.classList.remove('locked')}else{err.textContent='Access key did not match.';input.select()}})
})();
