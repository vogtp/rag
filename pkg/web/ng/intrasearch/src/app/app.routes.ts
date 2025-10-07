import { Routes } from '@angular/router';
import { CollectionOverviewComponent } from './components/collection-overview/collection-overview.component';
import { CollectionsListComponent } from './components/collections-list/collections-list.component';
import { SearchComponent } from './components/search/search.component';
import { SettingsComponent } from './components/settings/settings.component';

export const routes: Routes = [
  { path: 'search', component: CollectionsListComponent },
  { path: 'search/:collection', component: SearchComponent },
  { path: 'collection/:collection', component: CollectionOverviewComponent },
  { path: 'settings', component: SettingsComponent },
  { path: '', redirectTo: '/search', pathMatch: 'full' },
];
