import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIcon } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatToolbarModule } from '@angular/material/toolbar';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-settings-button',
  standalone: true,
  imports: [
    CommonModule,
    MatToolbarModule,
    MatIcon,
    MatMenuModule,
    MatButtonModule,
  ],
  templateUrl: './settings-button.component.html',
  styleUrl: './settings-button.component.css',
})
export class SettingsButtonComponent {
  constructor(private router: Router, private route: ActivatedRoute) {}

  // go() {
  //   this.router.navigate([`../edit/${item.id}`], { relativeTo: this.route });
  // }

  onSettingsClick() {
    console.log('Settings');
    this.router.navigate([`/settings`]);
  }

  onMenuItemClick(option: string) {
    console.log('Selected:', option);
  }
}
