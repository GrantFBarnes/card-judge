document.addEventListener("htmx:beforeOnLoad", function (event) {
  // always swap htmx response even if event.detail.xhr.status != 200
  event.detail.shouldSwap = true;
  event.detail.isError = false;
});
