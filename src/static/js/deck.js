document.addEventListener("htmx:afterSwap", function (event) {
	if (event.target.id === "deck-card-export-content") {
		downloadCSV("card-judge-cards.csv", event.target.innerHTML);
	}
});

function downloadCSV(fileName, content) {
	const element = document.createElement("a");
	element.setAttribute("href", "data:text/csv;charset=utf-8," + encodeURIComponent(content));
	element.setAttribute("download", fileName);
	element.style.display = "none";
	document.body.appendChild(element);
	element.click();
	document.body.removeChild(element);
}

function goToPage(pageNum) {
	const pageSearchEl = document.getElementById('pageSearch');
	if (pageSearchEl) {
		pageSearchEl.value = pageNum;
		htmx.trigger('#card-search-form', 'submit');
	}
}

function nextPage() {
	const pageJumpEl = document.getElementById('pageJump');
	if (pageJumpEl) {
		const currentPage = parseInt(pageJumpEl.value) || 1;
		const maxPage = parseInt(pageJumpEl.max) || 1;
		goToPage(Math.min(currentPage + 1, maxPage));
	}
}

function previousPage() {
	const pageJumpEl = document.getElementById('pageJump');
	if (pageJumpEl) {
		const currentPage = parseInt(pageJumpEl.value) || 1;
		goToPage(Math.max(currentPage - 1, 1));
	}
}