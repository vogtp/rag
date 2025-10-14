import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import {
  CollectionListResponse,
  CollectionListService,
} from '../../services/collection-list.service';

import { MatListModule } from '@angular/material/list';
import { RouterLink } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'app-collections-list',
  standalone: true,
  imports: [CommonModule, MatListModule, RouterLink],
  templateUrl: './collections-list.component.html',
  styleUrl: './collections-list.component.css',
})
export class CollectionsListComponent {
  collectionResponse: CollectionListResponse | undefined;

  constructor(private collectionService: CollectionListService) {}

  loadCollections() {
    this.collectionService.getCollections().subscribe({
      next: (data) => {
        this.collectionResponse = data;
      },
      error: (err) => {
        console.error('Load collection list from backend: ' + err);
        if (err === undefined) {
          return;
        }
        console.log('Load collection list Http status: ' + err.status);
        if (err instanceof HttpErrorResponse) {
          console.log('Load collection list Http status: ' + err.status);

          window.location.href = '/login?OrigPath=' + window.location.href;
        }
      },
    });
  }

  ngOnInit() {
    this.loadCollections();
  }
}
