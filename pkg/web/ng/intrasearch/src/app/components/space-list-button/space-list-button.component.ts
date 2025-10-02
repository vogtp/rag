import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIcon } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatToolbarModule } from '@angular/material/toolbar';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import {
  CollectionListResponse,
  CollectionListService,
} from '../../services/collection-list.service';

@Component({
  selector: 'app-space-list-button',
  standalone: true,
  imports: [
    CommonModule,
    RouterLink,
    MatToolbarModule,
    MatIcon,
    MatMenuModule,
    MatButtonModule,
  ],
  templateUrl: './space-list-button.component.html',
  styleUrl: './space-list-button.component.css',
})
export class SpaceListButtonComponent {
  collectionResponse: CollectionListResponse | undefined;

  constructor(private collectionService: CollectionListService) {}

  loadCollections() {
    this.collectionService.getCollections().subscribe({
      next: (data) => {
         this.collectionResponse = data;
      },
      error: (err) => {
        console.error(err);
        window.location.href = '/login?OrigPath=' + window.location.href;
      },
    });
  }

  ngOnInit() {
    this.loadCollections();
  }
}
