import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { httpHeaders, userSettingsURL } from './common';
import { User } from './user.structs';

@Injectable({
  providedIn: 'root',
})
export class SettingsService {
  getUserSetting(): Observable<User> {
    let url = userSettingsURL;
    return this.http.get<User>(url, { headers: httpHeaders });
  }

  saveUserSetting(us: User) {
    let url = userSettingsURL;
    // url = 'http://localhost:4444' + userSettingsURL;
    console.log('Sending usersettings save put: ' + url);
    return this.http.put(url, us, { headers: httpHeaders });
  }

  constructor(private http: HttpClient) {}
}
