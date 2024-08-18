require("../navbar")

document.addEventListener('DOMContentLoaded', (event) => {
    
    const shareLinks = document.querySelectorAll('.share-link');
    const copyShareLinkButtons = document.querySelectorAll('.btn-copy-share-link');
    const copyShareLinkEmptyIcons = document.querySelectorAll('.clipboard-empty');
    const copyShareLinkCheckIcons = document.querySelectorAll('.clipboard-check');
    
    shareLinks.forEach(shareLink => {
        shareLink.textContent = window.location.origin + "/" + shareLink.dataset.sharedListPath + "/" + shareLink.dataset.shareCode;
    });

    copyShareLinkButtons.forEach(copyShareLinkBtn => {
        copyShareLinkBtn.addEventListener('click', async (event) => {
            var result = writeClipboardText(window.location.origin + "/" + copyShareLinkBtn.dataset.sharedListPath + "/" + copyShareLinkBtn.dataset.shareCode);
            if (result) {
                copyShareLinkEmptyIcons.forEach(icon => icon.classList.add("hidden"));
                copyShareLinkCheckIcons.forEach(icon => icon.classList.remove("hidden"));
            }
        });
    });

    async function writeClipboardText(text) {
        try {
            await navigator.clipboard.writeText(text);
            return true;
        } catch (error) {
            console.error(error.message);
            return false;
        }
    }
});