import { ChangeDetectionStrategy, Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';

@Component({
  selector: 'app-settings',
  imports: [
    FormsModule,
    MatButtonModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatIconModule,
  ],
  changeDetection: ChangeDetectionStrategy.OnPush,
  templateUrl: './settings.component.html',
  styleUrl: './settings.component.css',
})
export class SettingsComponent {
  apiKey = 'nokey yet';
  confluenceApiKey = '';
  spacesStr = '';

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
