import { NgFor } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  inject,
  model,
} from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import {
  MatBottomSheet,
  MatBottomSheetModule,
  MatBottomSheetRef,
} from '@angular/material/bottom-sheet';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatTabsModule } from '@angular/material/tabs';
import { SbbLoadingIndicatorModule } from '@sbb-esta/angular/loading-indicator';
import { SettingsService } from '../../services/settings.service';
import { Collection, User } from '../../services/settings.service.structs';
import { ChatbotIcons } from '../chat/interfaces/library.interface';
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
    ReactiveFormsModule,
    CollectionComponent,
    NgFor,
    MatTabsModule,
    SbbLoadingIndicatorModule,
    MatBottomSheetModule,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './settings.component.html',
  styleUrl: './settings.component.css',
})
export class SettingsComponent {
  userSettings: User | undefined;
  collections = model<Collection[]>();
  waitingResponse: boolean = false;

  respMsg: string = '';

  icons: ChatbotIcons = {
    chatbotIcon: '../assets/icons/chatbot.svg',
    userIcon: '../assets/icons/user.svg',
  };
  basePath: string = 'http://localhost:4444/api/chat/completions';
  model: string = 'gpt-oss';

  constructor(
    private settingsService: SettingsService,
    private cdRef: ChangeDetectorRef
  ) {}

  loadSettings() {
    this.waitingResponse = true;
    this.settingsService.getUserSetting().subscribe({
      next: (data) => {
        this.userSettings = data;
        this.collections.set(data.edges.Collections!);
        this.cdRef.detectChanges();
      },
      error: (err) => {
        console.error('Load setting from backend: ' + err);
        this.waitingResponse = false;
        window.location.href = '/login?OrigPath=' + window.location.href;
      },
      complete: () => {
        console.debug('request usersettings complete');
        //window.location.reload();
        this.waitingResponse = false;
      },
    });
  }

  ngOnInit() {
    this.loadSettings();
  }
  private _bottomSheet = inject(MatBottomSheet);
  onSaveClick() {
    this._bottomSheet.open(SavingFeedbackBottomSheed);
    this.waitingResponse = true;
    this.respMsg = '';
    this.settingsService.saveUserSetting(this.userSettings!).subscribe({
      error: (err) => {
        if (err instanceof HttpErrorResponse) {
          this.respMsg = err.error + ' (' + err.statusText + ')';
        } else {
          this.respMsg = err;
        }
        console.error(err);
        this.waitingResponse = false;
      },
      complete: () => {
        this._bottomSheet.dismiss(SavingFeedbackBottomSheed);
        this._bottomSheet.open(SavedFeedbackBottomSheed);
        console.info('save complete');
        this.respMsg = 'Save success full';
        this.waitingResponse = false;
        // window.location.reload();
      },
    });
  }

  addCollection() {
    let col = new Collection();
    col.Name = 'New Collection (please change)';
    // let src = new SourceSystem();
    // src.Name = "New Source (please change)"
    // col.edges.Sources?.push(src)
    this.userSettings?.edges.Collections?.push(col);
  }

  debug() {
    console.log(this.userSettings);
  }
}

@Component({
  template: 'Saving Settings...',
  imports: [],
})
export class SavingFeedbackBottomSheed {}

@Component({
  template: '<div (click)="click()">Settings Saved</div>',
  imports: [],
})
export class SavedFeedbackBottomSheed {
  private _bottomSheetRef =
    inject<MatBottomSheetRef<SavedFeedbackBottomSheed>>(MatBottomSheetRef);

  click(): void {
    this._bottomSheetRef.dismiss();
  }
}
