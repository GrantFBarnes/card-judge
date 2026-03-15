function selectDeck(deckId) {
	const promptCheckbox = document.getElementById(`deckSelectPrompt${deckId}`);
	if (!promptCheckbox) return;
	const responseCheckbox = document.getElementById(`deckSelectResponse${deckId}`);
	if (!responseCheckbox) return;

	if (promptCheckbox.checked && responseCheckbox.checked) {
		promptCheckbox.checked = false;
		responseCheckbox.checked = false;
	} else {
		promptCheckbox.checked = true;
		responseCheckbox.checked = true;
	}
}
