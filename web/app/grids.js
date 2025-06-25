/**
 * Reusable grid functionality with master/detail pattern
 * Implements master/detail pattern for pre-rendered HTML tables
 */


/**
 * Initializes a pre-rendered HTML table with master/detail functionality
 * @param {string} gridSelector - CSS selector for the grid container
 * @param {Object} options - Configuration options
 * @param {string} [options.defaultSortColumn] - Column field to sort by default
 * @param {string} [options.defaultSortDirection] - Default sort direction ('asc' or 'desc')
 */
export function initMasterDetailGrid(gridSelector, options = {}) {
  const gridContainer = document.querySelector(gridSelector);
  if (!gridContainer) return null;
  
  // Find all tables in the grid container
  const table = gridContainer.querySelector('table');
  if (!table) return null;
  
  // Get table header and sortable columns
  const thead = table.querySelector('thead');
  const sortableHeaders = thead ? thead.querySelectorAll('th.sortable') : [];
  
  // Get table body and data rows
  const tbody = table.querySelector('tbody');
  const dataRows = tbody ? tbody.querySelectorAll('tr.data-row') : [];
  const detailRows = tbody ? tbody.querySelectorAll('tr.detail-row') : [];
  
  // Initialize state
  const state = {
    sortColumn: options.defaultSortColumn || null,
    sortDirection: options.defaultSortDirection || 'asc',
    expandedRows: new Set()
  };
  
  /**
   * Toggles the detail row visibility for a parent row
   * @param {string} rowId - The ID of the parent row
   */
  function toggleDetailRow(rowId) {
    const detailRow = tbody.querySelector(`tr.detail-row[data-parent-id="${rowId}"]`);
    if (!detailRow) return;
    
    if (state.expandedRows.has(rowId)) {
      // Collapse row
      state.expandedRows.delete(rowId);
      detailRow.classList.remove('detail-row-visible');
      detailRow.classList.add('detail-row-enter');
    } else {
      // Expand row
      state.expandedRows.add(rowId);
      detailRow.classList.remove('detail-row-enter');
      detailRow.classList.add('detail-row-visible');
    }
  }
  
  /**
   * Sort table rows based on a specific column
   * @param {string} field - The field to sort by
   * @param {string} direction - The sort direction ('asc' or 'desc')
   */
  function sortTableByColumn(field, direction) {
    const rowPairs = [];
    
    // Create pairs of [data-row, detail-row] for sorting
    dataRows.forEach(dataRow => {
      const rowId = dataRow.dataset.id;
      const detailRow = tbody.querySelector(`tr.detail-row[data-parent-id="${rowId}"]`);
      rowPairs.push([dataRow, detailRow]);
    });
    
    // Sort the row pairs
    rowPairs.sort((a, b) => {
      const rowA = a[0];
      const rowB = b[0];
      
      // Get values based on field
      let valueA, valueB;
      
      if (field === 'name') {
        // For name, get the text content of the first cell
        valueA = rowA.querySelector('td').textContent.trim().toLowerCase();
        valueB = rowB.querySelector('td').textContent.trim().toLowerCase();
        const result = valueA.localeCompare(valueB);
        return direction === 'asc' ? result : -result;
      } 
      else if (field === 'priority') {
        // For priority, use the priority-cell class to find the cell
        const priorityCellA = rowA.querySelector('.priority-cell');
        const priorityCellB = rowB.querySelector('.priority-cell');
        
        valueA = priorityCellA ? parseInt(priorityCellA.textContent.trim()) || Infinity : Infinity;
        valueB = priorityCellB ? parseInt(priorityCellB.textContent.trim()) || Infinity : Infinity;
        
        // Handle empty priorities (sort to bottom)
        const isEmptyA = isNaN(valueA) || valueA === Infinity;
        const isEmptyB = isNaN(valueB) || valueB === Infinity;
        
        if (isEmptyA && isEmptyB) return 0;
        if (isEmptyA) return 1; // A is empty, sort it after B
        if (isEmptyB) return -1; // B is empty, sort it after A
        
        return direction === 'asc' ? valueA - valueB : valueB - valueA;
      }
      
      return 0; // Default no sorting
    });
    
    // Remove all rows and reinsert them in the sorted order
    const fragment = document.createDocumentFragment();
    rowPairs.forEach(([dataRow, detailRow]) => {
      fragment.appendChild(dataRow);
      if (detailRow) {
        fragment.appendChild(detailRow);
      }
    });
    
    // Clear and append
    while (tbody.firstChild) {
      tbody.removeChild(tbody.firstChild);
    }
    tbody.appendChild(fragment);
  }
  
  /**
   * Updates sort indicators in column headers
   */
  function updateSortIndicators() {
    sortableHeaders.forEach(header => {
      const field = header.dataset.field;
      const indicator = header.querySelector('.sort-indicator');
      
      if (field === state.sortColumn) {
        // Show sort indicator
        indicator.classList.remove('opacity-0');
        
        // Update direction indicator
        indicator.innerHTML = state.sortDirection === 'asc' 
          ? `<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clip-rule="evenodd" />
            </svg>`
          : `<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>`;
        
        // Highlight header
        header.querySelector('div').classList.add('font-semibold');
      } else {
        // Hide sort indicator
        indicator.classList.add('opacity-0');
        header.querySelector('div').classList.remove('font-semibold');
      }
    });
  }
  
  // Add event listeners to data rows for toggle detail view
  dataRows.forEach(row => {
    row.addEventListener('click', (event) => {
      // Don't expand if clicking on action buttons
      if (event.target.closest('button') || event.target.closest('a')) {
        return;
      }
      const rowId = row.dataset.id;
      toggleDetailRow(rowId);
    });
  });
  
  // Add event listeners for sorting column headers
  sortableHeaders.forEach(header => {
    header.addEventListener('click', () => {
      const field = header.dataset.field;
      
      // Toggle direction if clicking the same column again
      if (state.sortColumn === field) {
        state.sortDirection = state.sortDirection === 'asc' ? 'desc' : 'asc';
      } else {
        state.sortColumn = field;
        state.sortDirection = 'asc';
      }
      
      // Sort the table
      sortTableByColumn(field, state.sortDirection);
      
      // Update indicators
      updateSortIndicators();
    });
  });
  
  // Apply initial sort if default is provided
  if (state.sortColumn) {
    sortTableByColumn(state.sortColumn, state.sortDirection);
    updateSortIndicators();
  }
  
  // Return public API
  return {
    toggleRow: (rowId) => toggleDetailRow(rowId),
    sort: (field, direction) => {
      state.sortColumn = field;
      state.sortDirection = direction || 'asc';
      sortTableByColumn(field, state.sortDirection);
      updateSortIndicators();
    },
    getExpandedRows: () => [...state.expandedRows],
    destroy: () => {
      // Clean up code if needed
    }
  };
}

// Helper function to format null/undefined values
export function formatNullableValue(value, defaultValue = '') {
  if (value === null || value === undefined) return defaultValue;
  if (typeof value === 'object' && 'Valid' in value) {
    return value.Valid ? (value.String || value.Int64 || '') : defaultValue;
  }
  return value;
}

// Common sort functions
export const sortFunctions = {
  // String sort
  string: (a, b, field) => {
    const valueA = formatNullableValue(a[field], '').toString().toLowerCase();
    const valueB = formatNullableValue(b[field], '').toString().toLowerCase();
    return valueA.localeCompare(valueB);
  },
  
  // Number sort
  number: (a, b, field) => {
    const valueA = parseFloat(formatNullableValue(a[field], 0));
    const valueB = parseFloat(formatNullableValue(b[field], 0));
    return valueA - valueB;
  },
  
  // Priority sort (handling nullable values)
  priority: (a, b) => {
    const p1 = a.priority && a.priority.Valid ? a.priority.Int64 : null;
    const p2 = b.priority && b.priority.Valid ? b.priority.Int64 : null;
    
    const p1Empty = p1 === undefined || p1 === null;
    const p2Empty = p2 === undefined || p2 === null;
    
    if (p1Empty && p2Empty) {
      return 0;
    } else if (p1Empty) {
      return 1; // place p1 after p2
    } else if (p2Empty) {
      return -1; // place p2 after p1
    } else {
      return p1 - p2;
    }
  }
};
