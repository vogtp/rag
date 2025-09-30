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
    return this.http.get<UserSettings>(url, {
      headers: httpHeaders,
    });
  }
  // summary(uuid: string): Observable<Document> {
  //   let url = summaryUrl + uuid;
  //   return this.http.get<Document>(url, { headers: httpHeaders });
  // }

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
  APIKey: string;
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