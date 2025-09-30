import { NgFor } from '@angular/common';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
} from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { SettingsService, UserSettings } from '../../services/settings.service';
import { CollectionComponent } from './collection/collection.component';
import { MatTabsModule } from '@angular/material/tabs';

@Component({
  selector: 'app-settings',
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
  apiKey = 'no key yet';
  confluenceURL = '';
  confluenceApiKey = '';
  spacesStr = '';

  userSettings: UserSettings | undefined;

  constructor(
    private settingsService: SettingsService,
    private cdRef: ChangeDetectorRef
  ) {}

  loadSettings() {
    this.settingsService.getUserSetting().subscribe((data) => {
      console.log(data);
      this.userSettings = data;
      this.confluenceURL = data.edges.Collections[0].edges.Sources[0].URL;
      this.cdRef.detectChanges();
    });
  }

  ngOnInit() {
    this.loadSettings();
  }

  onSaveClick() {
    console.log('API key ' + this.apiKey);
    console.log('Spaces ' + this.spacesStr);
  }

  onGenerateKey() {
    console.log('Generating API key');
  }
  onCopyKey() {
    console.log('Copy API key');
  }
}
