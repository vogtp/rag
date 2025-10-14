import { NgFor } from '@angular/common';
import { ChangeDetectionStrategy, Component, model } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatTabsModule } from '@angular/material/tabs';
import { MatTreeModule } from '@angular/material/tree';
import { Collection, SourceSystem } from '../../../services/user.structs';
import { SourceComponent } from '../source/source.component';

@Component({
  selector: 'app-collection',
  standalone: true,
  imports: [
    MatButtonModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatIconModule,
    FormsModule,
    ReactiveFormsModule,
    MatCardModule,
    MatTabsModule,
    MatTreeModule,
    MatIconModule,
    SourceComponent,
  ],
  templateUrl: './collection.component.html',
  styleUrl: './collection.component.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class CollectionComponent {
  collection = model<Collection>();

  // addSource() {
  //   console.log('Add source');
  //   let src = new SourceSystem();
  //   src.Name = 'New Source';
  //   this.collection()?.Source? = src;
  //   console.log(this.collection());
  // }
}
