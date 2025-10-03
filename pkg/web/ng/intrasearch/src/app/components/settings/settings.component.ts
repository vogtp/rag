import { NgFor } from '@angular/common';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  model,
} from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatTabsModule } from '@angular/material/tabs';
import { SettingsService } from '../../services/settings.service';
import { Collection, User } from '../../services/settings.service.structs';
import { CollectionComponent } from './collection/collection.component';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [
    FormsModule,
    MatButtonModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatIconModule,
    FormsModule,
    ReactiveFormsModule,
    CollectionComponent,
    NgFor,
    MatTabsModule,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './settings.component.html',
  styleUrl: './settings.component.css',
})
export class SettingsComponent {
  userSettings: User | undefined;
  collections = model<Collection[]>();
  //userSettings = model<UserSettings>();

  constructor(
    private settingsService: SettingsService,
    private cdRef: ChangeDetectorRef
  ) {}

  loadSettings() {
    this.settingsService.getUserSetting().subscribe({
      next: (data) => {
        this.userSettings = data;
        this.collections.set(data.edges.Collections!);
        this.cdRef.detectChanges();
      },
      error: (err) => {
        console.error(err);
        window.location.href = '/login?OrigPath=' + window.location.href;
      },
      complete: () => console.debug('request usersettings complete'),
    });
  }

  ngOnInit() {
    this.loadSettings();
  }

  onSaveClick() {
    this.settingsService.saveUserSetting(this.userSettings!);
  }

  debug() {
    console.log(this.userSettings);
  }
}
