import { NgFor, NgIf } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  Inject,
  inject,
  model,
} from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import {
  MAT_BOTTOM_SHEET_DATA,
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
import { Collection, SourceSystem, User } from '../../services/user.structs';
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
    NgIf
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

  constructor(
    private settingsService: SettingsService,
    private cdRef: ChangeDetectorRef
  ) {}

  loadSettings() {
    this.waitingResponse = true;
    this.settingsService.getUserSetting().subscribe({
      next: (data) => {
        this.userSettings = data;
        this.collections.set(data.Collections!);
        this.cdRef.detectChanges();
      },
      error: (err) => {
        console.log('Load setting from backend: ' + err);
        this.waitingResponse = false;
        window.location.href = '/login?OrigPath=' + window.location.href;
      },
      complete: () => {
        console.log('request usersettings complete');
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
    this._bottomSheet.open(UserMsgBottomSheed, {
      data: { message: 'Saving settings' },
    });
    this.waitingResponse = true;
    this.respMsg = '';
    this.settingsService.saveUserSetting(this.userSettings!).subscribe({
      error: (err) => {
        if (err instanceof HttpErrorResponse) {
          this.respMsg = err.error + ' (' + err.statusText + ')';
        } else {
          this.respMsg = err;
        }
        this._bottomSheet.open(UserMsgBottomSheed, {
          data: { message: this.respMsg },
        });
        console.error(err);
        this.waitingResponse = false;
      },
      complete: () => {
        this._bottomSheet.open(UserMsgBottomSheed, {
          data: { message: 'Saved settings' },
        });
        console.info('save complete');
        this.respMsg = 'Save success full';
        this.waitingResponse = false;
        // window.location.reload();
      },
    });
  }

  addCollection() {
    let col = new Collection();
    col.Displayname = 'New Collection (please change)';
    let src = new SourceSystem();
    src.Name = 'New Source (please change)';
    let s = this.userSettings!.Collections![0].Source!;
    src.URL = s.URL;
    src.Type = s.Type;
    src.Name = s.Name;
    col.Source = src;
    console.log(col);
    // let src = new SourceSystem();
    // src.Name = "New Source (please change)"
    // col.edges.Sources?.push(src)
    this.userSettings?.Collections?.unshift(col);
  }

  debug() {
    console.log(this.userSettings);
  }
}

@Component({
  template: '<div (click)="click($event)">{{data.message}}</div>',
  imports: [],
})
export class UserMsgBottomSheed {
  constructor(
    @Inject(MAT_BOTTOM_SHEET_DATA) public data: { message: string }
  ) {}
  private _bottomSheetRef =
    inject<MatBottomSheetRef<UserMsgBottomSheed>>(MatBottomSheetRef);

  click(event: MouseEvent): void {
    this._bottomSheetRef.dismiss();
    event.preventDefault();
  }
}
