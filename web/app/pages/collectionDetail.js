require("../navbar");

document.addEventListener('DOMContentLoaded', () => {
  const shareLinks = document.querySelectorAll('.share-link');
  shareLinks.forEach(shareLink => {
    shareLink.textContent = window.location.origin + "/" + shareLink.dataset.sharedListPath + "/" + shareLink.dataset.shareCode;
  });
});