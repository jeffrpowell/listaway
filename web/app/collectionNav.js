// Add collection links to navigation
document.addEventListener('DOMContentLoaded', () => {
  // Find the navigation element - typically after the Lists link
  const navContainer = document.querySelector('.navbar-nav, nav ul, .nav-links');
  if (!navContainer) return;

  // Find the Lists link to position our Collections link after it
  const listItems = Array.from(navContainer.querySelectorAll('li, .nav-item'));
  const listsNavItem = listItems.find(item => {
    const link = item.querySelector('a');
    return link && (link.textContent.includes('Lists') || link.href.includes('/lists'));
  });

  if (listsNavItem) {
    // Create the Collections nav item
    const collectionsNavItem = document.createElement(listsNavItem.tagName);
    collectionsNavItem.className = listsNavItem.className;
    collectionsNavItem.innerHTML = `
      <a href="/collections" class="${listsNavItem.querySelector('a').className}">
        Collections
      </a>
    `;

    // Insert after the Lists nav item
    if (listsNavItem.nextSibling) {
      navContainer.insertBefore(collectionsNavItem, listsNavItem.nextSibling);
    } else {
      navContainer.appendChild(collectionsNavItem);
    }
  }
});
