require('../index')
require('../navbar')
require('../grids')
import { initMasterDetailGrid } from "../grids";

document.addEventListener('DOMContentLoaded', (event) => {    
    // Initialize the master/detail grid
    const gridApi = initMasterDetailGrid('.item-grid', {
        defaultSortColumn: 'priority',
        defaultSortDirection: 'asc'
    });
});