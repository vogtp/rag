import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { httpHeaders, userSettingsURL } from './common';

@Injectable({
  providedIn: 'root',
})
export class SettingsService {
  getUserSetting(): Observable<UserSettings> {
    let url = userSettingsURL;
    return this.http.get<UserSettings>(url, { headers: httpHeaders });
  }

  saveUserSetting(us: UserSettings) {
    let url = userSettingsURL;
    // url = 'http://localhost:4444' + userSettingsURL;
    console.log('Sending usersettings save put: ' + url);
    this.http.put(url, us, { headers: httpHeaders }).subscribe({
      next: (v) => console.log(v),
      error: (e) => console.error(e),
      complete: () => console.info('complete'),
    });
  }

  constructor(private http: HttpClient) {}
}

export interface UserSettings {
  id: number;
  Name: string;
  edges: Edges;
}

export interface Edges {
  Collections: Collection[];
}

export interface Collection {
  id: number;
  Name: string;
  edges: EdgesCol;
}

export interface EdgesCol {
  Sources: Source[];
}

export interface Source {
  id: number;
  Name: string;
  Type: string;
  URL: string;
  key?: string;
  parts: string;
}
