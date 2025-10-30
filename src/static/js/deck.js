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

function goToPage(pageNumber) {
	const pageSearchElement = document.getElementById("pageSearch");
	if (pageSearchElement) {
		if (pageSearchElement.value != pageNumber) {
			pageSearchElement.value = pageNumber;
			htmx.trigger("#card-search-form", "submit");
		}
	}
}

function previousPage() {
	const pageSearchElement = document.getElementById("pageSearch");
	if (pageSearchElement) {
		const currentPage = parseInt(pageSearchElement.value) || 1;
		goToPage(Math.max(currentPage - 1, 1));
	}
}

function nextPage() {
	const pageSearchElement = document.getElementById("pageSearch");
	if (pageSearchElement) {
		const currentPage = parseInt(pageSearchElement.value) || 1;
		const totalPages = parseInt(pageSearchElement.max) || 1;
		goToPage(Math.min(currentPage + 1, totalPages));
	}
}