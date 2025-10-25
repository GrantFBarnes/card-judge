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

// Pagination functions
let currentPageNumber = 1;
let currentPageSize = 10;

function updatePageDisplay() {
    const currentPageEl = document.getElementById('currentPage');
    const pageJumpEl = document.getElementById('pageJump');
    const pageSizeSelectEl = document.getElementById('pageSizeSelect');
    
    if (currentPageEl) currentPageEl.textContent = currentPageNumber;
    if (pageJumpEl) pageJumpEl.value = currentPageNumber;
    if (pageSizeSelectEl) pageSizeSelectEl.value = currentPageSize;
}

function goToPage(pageNum) {
    if (pageNum < 1) pageNum = 1;
    if (pageNum > 100) pageNum = 100;
    
    currentPageNumber = pageNum;
    const pageSearchEl = document.getElementById('pageSearch');
    const pageSizeSearchEl = document.getElementById('pageSizeSearch');
    
    if (pageSearchEl) {
        pageSearchEl.value = currentPageNumber;
    }
    if (pageSizeSearchEl) {
        pageSizeSearchEl.value = currentPageSize;
    }
    
    updatePageDisplay();
    
    // Trigger the htmx form submission
    const formEl = document.getElementById('card-search-form');
    if (formEl) {
        htmx.trigger(formEl, 'submit');
    }
}

function nextPage() {
    goToPage(currentPageNumber + 1);
}

function previousPage() {
    goToPage(currentPageNumber - 1);
}

function changePageSize(newSize) {
    if (newSize < 1) newSize = 10;
    if (newSize > 50) newSize = 50;
    
    currentPageSize = newSize;
    currentPageNumber = 1; // Reset to first page when changing page size
    
    const pageSearchEl = document.getElementById('pageSearch');
    const pageSizeSearchEl = document.getElementById('pageSizeSearch');
    
    if (pageSearchEl) {
        pageSearchEl.value = 1;
    }
    if (pageSizeSearchEl) {
        pageSizeSearchEl.value = currentPageSize;
    }
    
    updatePageDisplay();
    
    // Trigger the htmx form submission
    const formEl = document.getElementById('card-search-form');
    if (formEl) {
        htmx.trigger(formEl, 'submit');
    }
}

// Reset to page 1 when search criteria changes
document.addEventListener('DOMContentLoaded', function() {
    const categorySearch = document.getElementById('categorySearch');
    const textSearch = document.getElementById('textSearch');
    
    if (categorySearch) {
        categorySearch.addEventListener('change', function() {
            currentPageNumber = 1;
            const pageSearchEl = document.getElementById('pageSearch');
            if (pageSearchEl) pageSearchEl.value = 1;
            updatePageDisplay();
        });
    }
    
    if (textSearch) {
        let searchTimeout;
        textSearch.addEventListener('input', function() {
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(function() {
                currentPageNumber = 1;
                const pageSearchEl = document.getElementById('pageSearch');
                if (pageSearchEl) pageSearchEl.value = 1;
                updatePageDisplay();
            }, 500);
        });
    }
});
