import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { MatIcon } from '@angular/material/icon';
import { ActivatedRoute, Router } from '@angular/router';
import {
  CollectionSearchResponse,
  CollectionSearchService,
} from '../../services/collection-search.service';
import { SearchResultItemComponent } from '../search-result-item/search-result-item.component';

import { HttpErrorResponse } from '@angular/common/http';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { Observable, catchError, throwError } from 'rxjs';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [
    CommonModule,
    MatIcon,
    ReactiveFormsModule,
    SearchResultItemComponent,
    MatFormFieldModule,
    MatIcon,
    MatInputModule,
  ],
  templateUrl: './search.component.html',
  styleUrl: './search.component.css',
})
export class SearchComponent {
  @Input() collection: string = 'intranet-all';
  @Input()
  set query(q: string) {
    this.searchQuery.setValue(q);
    if (q) {
      console.log('Searching for ' + q);

      this.search();
    }
  }
  searchQuery = new FormControl('');
  searchResult: CollectionSearchResponse | undefined;
  message: string = '';

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private searchService: CollectionSearchService
  ) {
    route.params.subscribe({
      next: (data) => {
        this.collection = this.route.snapshot.params['collection'];
        this.search();
      },
      error: (err) => {
        console.error('Load search from backend: ' + err);
        window.location.href = '/login?OrigPath=' + window.location.href;
      },
      complete: () => console.debug('search request complete'),
    });
  }

  search() {
    let query = this.searchQuery.value!;
    console.log('Query: ' + query);
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: {
        query: query,
      },
      queryParamsHandling: 'merge',
      // preserve the existing query params in the route
      skipLocationChange: false,
      // do not trigger navigation
    });
    this.searchResult = undefined;
    this.message = '';
    this.searchService
      .searchCollection(this.collection, query)
      .pipe(catchError(this.handleError))
      .subscribe((data) => {
        this.searchResult = data;
        if (data && data.Documents) {
          this.message = 'Results: ' + data.Documents?.length;
        }
      });
  }

  private handleError(
    error: any,
    caught: Observable<CollectionSearchResponse>
  ) {
    let err = '';
    if (error.error instanceof ErrorEvent) {
      err = error.error.message;
    } else if (error instanceof HttpErrorResponse) {
      err = error.error.Error;
    } else {
      err = error.status;
    }
    console.log('err: ' + err);

    console.log('message: ' + this.message);
    this.message = err;
    console.log('message: ' + this.message);

    return throwError(() => new Error('test'));
  }
}
