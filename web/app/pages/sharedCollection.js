require("../navbar");

// Store list data
let listsWithItems = new Map();

// Event handler when the page loads
document.addEventListener('DOMContentLoaded', () => {
  // Load items for each list in the collection
  loadAllListItems();
});

// Load items for all lists in the collection
function loadAllListItems() {
  const listContainers = document.querySelectorAll('[id^="listItems-"]');
  
  listContainers.forEach(container => {
    const listId = container.id.split('-')[1];
    if (listId) {
      fetchListItems(listId, container);
    }
  });
}

// Fetch items for a specific list
function fetchListItems(listId, container) {
  // First check if this list has a share code
  const listElement = container.closest('.bg-white');
  const shareLink = listElement ? listElement.querySelector('a[href^="/sharedlist/"]') : null;
  
  if (shareLink) {
    // Extract share code from the link
    const shareCode = shareLink.getAttribute('href').split('/').pop();
    
    // Fetch items using the share code
    fetch(`/sharedlist/${shareCode}/items`)
      .then(response => {
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        return response.json();
      })
      .then(items => {
        listsWithItems.set(parseInt(listId), items);
        renderListItems(container, items);
      })
      .catch(error => {
        console.error(`Error loading items for list ${listId}:`, error);
        container.innerHTML = `<p class="text-center text-red-500 py-4">Error loading items. Please try again.</p>`;
      });
  } else {
    container.innerHTML = `<p class="text-center text-font-secondary-light py-4">This list is not publicly shared.</p>`;
  }
}

// Render items for a list
function renderListItems(container, items) {
  if (!items || items.length === 0) {
    container.innerHTML = `<p class="text-center text-font-secondary-light py-4">This list has no items.</p>`;
    return;
  }
  
  // Sort items by priority if available
  const sortedItems = [...items].sort((a, b) => {
    const priorityA = a.priority && a.priority.Valid ? a.priority.Int64 : 999;
    const priorityB = b.priority && b.priority.Valid ? b.priority.Int64 : 999;
    return priorityA - priorityB;
  });
  
  let html = '<ul class="divide-y divide-gray-200">';
  
  sortedItems.forEach(item => {
    const hasPriority = item.priority && item.priority.Valid;
    const hasUrl = item.url && item.url.Valid;
    const hasNotes = item.notes && item.notes.Valid;
    
    html += `
      <li class="py-4">
        <div class="flex items-start">
          ${hasPriority ? `<span class="px-2 py-1 text-xs font-semibold bg-blue-100 text-font-link rounded-full mr-2">${item.priority.Int64}</span>` : ''}
          <div class="flex-1">
            <h3 class="text-lg font-medium">
              ${hasUrl 
                ? `<a href="${item.url.String}" target="_blank" class="text-font-link hover:text-font-link">${item.name} <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 inline-block" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" /></svg></a>` 
                : item.name
              }
            </h3>
            ${hasNotes ? `<p class="text-font-secondary-light mt-1">${item.notes.String}</p>` : ''}
          </div>
        </div>
      </li>
    `;
  });
  
  html += '</ul>';
  container.innerHTML = html;
}
