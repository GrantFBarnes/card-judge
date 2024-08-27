function toggleDarkMode() {
  const className = "dark-theme";

  document.body.classList.toggle(className);

  const toggleElement = document.getElementById("topbar-toggle-dark-mode");
  toggleElement.innerHTML = document.body.classList.contains(className)
    ? "&#127762;"
    : "&#127766;";
}
