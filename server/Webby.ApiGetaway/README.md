
# Adding Routes in API Gateway with Ocelot  
  
This guide explains how to add new routes to the API Gateway using **Ocelot**, including  jwt and role-based access control.  
  
---  
  
## 1. Route Definition  
  
Routes are defined in the `ocelot.json` configuration file. Each route typically includes:  
  
- **DownstreamPathTemplate**: The internal service endpoint  
- **DownstreamScheme**: `http` or `https`  
- **DownstreamHostAndPorts**: Host and port of the downstream service  
- **UpstreamPathTemplate**: The URL exposed via the API Gateway  
- **UpstreamHttpMethod**: Allowed HTTP methods (`GET`, `POST`, etc.)  
- **AuthenticationOptions** (optional): Specify if JWT authentication is required  
- **RouteClaimsRequirement** (optional): Define role-based access  
---  
  
## 2. Example: Protected Route for Moderators  
  
```json id="9afq1x"  
{  
   "DownstreamPathTemplate": "/api/auth/check",  
   "DownstreamScheme": "http",  
   "DownstreamHostAndPorts": [  
      {  
         "Host": "localhost",  
         "Port": 5001  
  }  
   ],  
   "UpstreamPathTemplate": "/api/auth/check",  
   "UpstreamHttpMethod": [  
      "GET"  
  ],  
   "AuthenticationOptions": {  
      "AuthenticationProviderKey": "Bearer"  
  },  
   "RouteClaimsRequirement": {  
      "Role": "Moderator"  
  }  
}
```

**Behavior:**

-   Users with `"Role": "Moderator"` → access granted
    
-   Users with `"Role": "Admin"` → access granted (inherits moderator rights)
    
-   Users with `"Role": "User"` → access denied (403)

## 3. Example: Protected Route for Admins Only

```json id="9afq1x"  
{  
   "DownstreamPathTemplate": "/api/admin/dashboard",  
   "DownstreamScheme": "http",  
   "DownstreamHostAndPorts": [  
      {  
         "Host": "localhost",  
         "Port": 5004 
  }  
   ],  
   "UpstreamPathTemplate": "/api/admin/dashboard",  
   "UpstreamHttpMethod": [  
      "GET"  
  ],  
   "AuthenticationOptions": {  
      "AuthenticationProviderKey": "Bearer"  
  },  
   "RouteClaimsRequirement": {  
      "Role": "Admin"  
  }  
}
```
**Behavior:**

-   Only users with `"Role": "Admin"` → access granted
    
-   Moderator or User → access denied (403)

## 4. Example: Open Route (No Authentication)

```json id="9afq1x"  
{  
   "DownstreamPathTemplate": "/api/public/info",  
   "DownstreamScheme": "http",  
   "DownstreamHostAndPorts": [  
      {  
         "Host": "localhost",  
         "Port": 5005
  }  
   ],  
   "UpstreamPathTemplate": "/api/public/info",  
   "UpstreamHttpMethod": [  
      "GET"  
  ] 
}
```
**Behavior:**
-   Accessible by anyone without a token
    
-   No `AuthenticationOptions` or `RouteClaimsRequirement` needed