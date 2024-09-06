document.addEventListener("htmx:beforeOnLoad", function (event) {
  // always swap htmx response even if event.detail.xhr.status != 200
  event.detail.shouldSwap = true;
  event.detail.isError = false;
});

document.addEventListener("htmx:afterSwap", function (event) {
  if (event?.detail?.pathInfo?.requestPath?.indexOf("/search") >= 0) return;
  if (event.detail.xhr.status >= 200 && event.detail.xhr.status < 300) {
    addClassToTarget("htmx-result-good", event.detail.target);
    removeClassFromTarget("htmx-result-bad", event.detail.target);
  } else {
    addClassToTarget("htmx-result-bad", event.detail.target);
    removeClassFromTarget("htmx-result-good", event.detail.target);
  }
});

function addClassToTarget(className, targetElement) {
  if (!targetElement.classList.contains(className)) {
    targetElement.classList.add(className);
  }
}

function removeClassFromTarget(className, targetElement) {
  if (targetElement.classList.contains(className)) {
    targetElement.classList.remove(className);
  }
}
