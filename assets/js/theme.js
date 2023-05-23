const themeSwitchers = document.querySelectorAll('.changeTheme');
let amogusLogo = document.querySelector(".amogusLogo");
let amogusLogo_webp = document.querySelector(".amogusLogo-webp");

themeSwitchers.forEach(switcher => {
  switcher.addEventListener('click', function() {
    applyTheme(this.dataset.theme);
    localStorage.setItem('theme', this.dataset.theme)
  });
});

function applyTheme(themeName) {
  let themeUrl = `../../../../assets/css/${themeName}.css`
  document.querySelector('[title="theme"]').setAttribute('href', themeUrl)
  amogusLogo.src= `../assets/img/AmogusLogo-${themeName}.jpg`;
  amogusLogo_webp.srcset= `../assets/img/AmogusLogo-${themeName}.webp`;
};

let activeTheme = localStorage.getItem('theme');

if (activeTheme === null) {
  applyTheme('light');
} else {
  applyTheme(activeTheme)
}