(() => {
  const q = document.querySelector('#search');
  const fam = document.querySelector('#family-filter');
  const stat = document.querySelector('#status-filter');
  const retired = document.querySelector('#show-retired');
  const cards = [...document.querySelectorAll('.project-card')];
  if (!cards.length) return;
  function apply(){
    const needle=(q?.value||'').trim().toLowerCase();
    for(const c of cards){
      const okQ=!needle || (c.dataset.search||'').toLowerCase().includes(needle);
      const okF=!fam?.value || c.dataset.family===fam.value;
      const okS=!stat?.value || c.dataset.status===stat.value;
      const okR=retired?.checked!==false || !['retired','archived'].includes(c.dataset.status);
      c.hidden=!(okQ&&okF&&okS&&okR);
    }
  }
  [q,fam,stat,retired].forEach(el=>el&&el.addEventListener('input',apply));
})();
